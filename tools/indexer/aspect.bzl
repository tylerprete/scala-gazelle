load("@bazel_tools//tools/cpp:toolchain_utils.bzl", "find_cpp_toolchain")
load(
    ":artifacts.bzl",
    "artifact_location",
    "sources_from_target",
    "struct_omit_none",
)

JarIndexerAspect = provider(
    "a provider for Jar Indexing",
    fields = {
        "info_file": "The index database",
        "jar_index_files": "A list of jar index files",
    },
)

_cpp_header_extensions = [
    "hh",
    "hxx",
    "ipp",
    "hpp",
]

_c_or_cpp_header_extensions = ["h"] + _cpp_header_extensions

_cpp_extensions = [
    "cc",
    "cpp",
    "cxx",
] + _cpp_header_extensions

_java_rules = [
    "java_library",
    "java_binary",
    "java_test",
    "java_proto_library",
    "jvm_import",
]

_scala_rules = [
    "scala_library",
    "scala_binary",
]

def get_aspect_ids(ctx, target):
    """Returns the all aspect ids, filtering out self."""
    aspect_ids = None
    if hasattr(ctx, "aspect_ids"):
        aspect_ids = ctx.aspect_ids
    elif hasattr(target, "aspect_ids"):
        aspect_ids = target.aspect_ids
    else:
        return None
    return [aspect_id for aspect_id in aspect_ids if "intellij_info_aspect" not in aspect_id]

def make_dep(dep, dependency_type):
    """Returns a Dependency proto struct."""
    return struct(
        dependency_type = dependency_type,
        target = dep.intellij_info.target_key,
    )

def make_deps(deps, dependency_type):
    """Returns a list of Dependency proto structs."""
    return [make_dep(dep, dependency_type) for dep in deps]

# Run-time dependency attributes, grouped by type.
RUNTIME_DEPS = [
    "runtime_deps",
]

PREREQUISITE_DEPS = []

# Dependency type enum
COMPILE_TIME = 0

RUNTIME = 1

# Compile-time dependency attributes, grouped by type.
DEPS = [
    "_cc_toolchain",  # From cc rules
    "_stl",  # From cc rules
    "malloc",  # From cc_binary rules
    "_java_toolchain",  # From java rules
    "deps",
    "jars",  # from java_import rules
    "exports",
    "java_lib",  # From old proto_library rules
    "_android_sdk",  # from android rules
    "aidl_lib",  # from android_sdk
    "_scala_toolchain",  # From scala rules
    "test_app",  # android_instrumentation_test
    "instruments",  # android_instrumentation_test
    "tests",  # From test_suite
]

# Defensive list of features that can appear in the C++ toolchain, but which we
# definitely don't want to enable (when enabled, they'd contribute command line
# flags that don't make sense in the context of intellij info).
UNSUPPORTED_FEATURES = [
    "thin_lto",
    "module_maps",
    "use_header_modules",
    "fdo_instrument",
    "fdo_optimize",
]

_all_rules = _java_rules + _scala_rules

def make_target_key(label, aspect_ids):
    """Returns a TargetKey proto struct from a target."""
    return struct_omit_none(
        aspect_ids = tuple(aspect_ids) if aspect_ids else None,
        label = str(label),
    )

def library_artifact(java_output):
    """Creates a LibraryArtifact representing a given java_output."""
    if java_output == None or java_output.class_jar == None:
        return None
    src_jars = get_source_jars(java_output)
    return struct_omit_none(
        interface_jar = artifact_location(java_output.ijar),
        jar = artifact_location(java_output.class_jar),
        source_jar = artifact_location(src_jars[0]) if src_jars else None,
        source_jars = [artifact_location(f) for f in src_jars],
    )

def semantics_extra_deps(base, semantics, name):
    if not hasattr(semantics, name):
        return base
    extra_deps = getattr(semantics, name)
    return base + extra_deps

def _is_proto_library_wrapper(target, ctx):
    """Returns True if the target is an empty shim around a proto library."""
    if not ctx.rule.kind.endswith("proto_library") or ctx.rule.kind == "proto_library":
        return False

    # treat any *proto_library rule with a single proto_library dep as a shim
    deps = collect_targets_from_attrs(ctx.rule.attr, ["deps"])
    return len(deps) == 1 and deps[0].intellij_info and deps[0].intellij_info.kind == "proto_library"

def _get_forwarded_deps(target, ctx):
    """Returns the list of deps of this target to forward.

    Used to handle wrapper/shim targets which are really just pointers to a
    different target (for example, java_proto_library)
    """
    if _is_proto_library_wrapper(target, ctx):
        return collect_targets_from_attrs(ctx.rule.attr, ["deps"])
    return []

def get_source_jars(output):
    """Returns a list of source jars from the output."""
    if hasattr(output, "source_jars"):
        return output.source_jars
    if hasattr(output, "source_jar"):
        return [output.source_jar]
    return []

def jars_from_output(output):
    """Collect jars for intellij-resolve-files from Java output."""
    if output == None:
        return []
    return [
        jar
        for jar in ([output.class_jar, output.ijar] + get_source_jars(output))
        if jar != None and not jar.is_source
    ]

def _collect_generated_files(java):
    """Collects generated files from a Java target"""
    if hasattr(java, "java_outputs"):
        return [
            (outputs.generated_class_jar, outputs.generated_source_jar)
            for outputs in java.java_outputs
            if outputs.generated_class_jar != None
        ]

    # Handles Bazel versions before 5.0.0.
    if (hasattr(java, "annotation_processing") and java.annotation_processing and java.annotation_processing.enabled):
        return [(java.annotation_processing.class_jar, java.annotation_processing.source_jar)]
    return []

def annotation_processing_jars(generated_class_jar, generated_source_jar):
    """Creates a LibraryArtifact representing Java annotation processing jars."""
    src_jar = generated_source_jar
    return struct_omit_none(
        jar = artifact_location(generated_class_jar),
        source_jar = artifact_location(src_jar),
        source_jars = [artifact_location(src_jar)] if src_jar else None,
    )

def get_java_provider(target):
    """Find a provider exposing java compilation/outputs data."""

    # Check for scala and kt providers before JavaInfo. e.g. scala targets have
    # JavaInfo, but their data lives in the "scala" provider and not JavaInfo.
    # See https://github.com/bazelbuild/intellij/pull/1202
    if hasattr(target, "scala"):
        return target.scala
    if hasattr(target, "kt") and hasattr(target.kt, "outputs"):
        return target.kt
    if JavaInfo in target:
        return target[JavaInfo]
    if hasattr(java_common, "JavaPluginInfo") and java_common.JavaPluginInfo in target:
        return target[java_common.JavaPluginInfo]
    return None

def update_set_in_dict(input_dict, key, other_set):
    """Updates depset in dict, merging it with another depset."""
    input_dict[key] = depset(transitive = [input_dict.get(key, depset()), other_set])

def divide_java_sources(ctx):
    """Divide sources into plain java, generated java, and srcjars."""

    java_sources = []
    gen_java_sources = []
    srcjars = []
    if hasattr(ctx.rule.attr, "srcs"):
        srcs = ctx.rule.attr.srcs
        for src in srcs:
            for f in src.files.to_list():
                if f.basename.endswith(".java"):
                    if f.is_source:
                        java_sources.append(f)
                    else:
                        gen_java_sources.append(f)
                elif f.basename.endswith(".srcjar"):
                    srcjars.append(f)

    return java_sources, gen_java_sources, srcjars

def _is_cpp_target(srcs):
    if all([src.extension in _c_or_cpp_header_extensions for src in srcs]):
        return True  # assume header-only lib is c++
    return any([src.extension in _cpp_extensions for src in srcs])

def _is_objcpp_target(srcs):
    return any([src.extension == "mm" for src in srcs])

def _sources(ctx, target):
    srcs = []
    if hasattr(ctx.rule.attr, "srcs"):
        srcs += [f for src in ctx.rule.attr.srcs for f in src.files.to_list()]
    if hasattr(ctx.rule.attr, "hdrs"):
        srcs += [f for src in ctx.rule.attr.hdrs for f in src.files.to_list()]

    return srcs

def is_valid_aspect_target(target):
    """Returns whether the target has had the aspect run on it."""
    return hasattr(target, "intellij_info")

def _collect_target_from_attr(rule_attrs, attr_name, result):
    """Collects the targets from the given attr into the result."""
    if not hasattr(rule_attrs, attr_name):
        return
    attr_value = getattr(rule_attrs, attr_name)
    type_name = type(attr_value)
    if type_name == "Target":
        result.append(attr_value)
    elif type_name == "list":
        result.extend(attr_value)

def collect_targets_from_attrs(rule_attrs, attrs):
    """Returns a list of targets from the given attributes."""
    result = []
    for attr_name in attrs:
        _collect_target_from_attr(rule_attrs, attr_name, result)
    return [target for target in result if is_valid_aspect_target(target)]

# Function copied from https://gist.github.com/oquenchil/7e2c2bd761aa1341b458cc25608da50c
# TODO: Directly use create_compile_variables and get_memory_inefficient_command_line.
def _get_compile_flags(dep):
    options = []
    compilation_context = dep[CcInfo].compilation_context
    for define in compilation_context.defines.to_list():
        options.append("-D\"{}\"".format(define))

    for define in compilation_context.local_defines.to_list():
        options.append("-D\"{}\"".format(define))

    for system_include in compilation_context.system_includes.to_list():
        if len(system_include) == 0:
            system_include = "."
        options.append("-isystem {}".format(system_include))

    for include in compilation_context.includes.to_list():
        if len(include) == 0:
            include = "."
        options.append("-I {}".format(include))

    for quote_include in compilation_context.quote_includes.to_list():
        if len(quote_include) == 0:
            quote_include = "."
        options.append("-iquote {}".format(quote_include))

    for framework_include in compilation_context.framework_includes.to_list():
        options.append("-F\"{}\"".format(framework_include))

    return options

def build_jar_index(ctx, target, jar):
    """Builds the java package manifest for the given source files."""

    ijar = ctx.actions.declare_file(jar.short_path.replace("/", "-"))

    # input_file = ctx.actions.declare_file(target.label.name + ".jar.json")
    output_file = ctx.actions.declare_file(target.label.name + ".jarindex.json")

    ctx.actions.run(
        executable = ctx.executable._ijar,
        inputs = [jar],
        outputs = [ijar],
        arguments = [
            "--target_label",
            str(ctx.label),
            jar.path,
            ijar.path,
        ],
        mnemonic = "Ijar",
    )

    # ctx.actions.write(
    #     content = json.encode(struct(
    #         label = str(target.label),
    #         filename = ijar_file.path,
    #     )),
    #     output = input_file,
    # )

    # args = [
    #     "--input_file",
    #     input_file.path,
    #     "--output_file",
    #     output_file.path,
    # ]

    # ctx.actions.run(
    #     mnemonic = "IJarIndexer",
    #     progress_message = "Extracting symbols from: " + jar.basename,
    #     executable = ctx.executable._jarindexer,
    #     arguments = args,
    #     inputs = [input_file, ijar_file],
    #     outputs = [output_file],
    # )

    ctx.actions.run(
        mnemonic = "JarIndexer2",
        progress_message = "Indexing " + ijar.basename,
        executable = ctx.executable._jarindexer2,
        arguments = [
            "--label",
            str(ctx.label),
            "--output_file",
            output_file.path,
            ijar.path,
        ],
        inputs = [ijar],
        outputs = [output_file],
    )

    return output_file

def collect_java_toolchain_info(target, ide_info, ide_info_file):
    """Updates java_toolchain-relevant output groups, returns false if not a java_toolchain target."""
    if hasattr(target, "java_toolchain"):
        toolchain = target.java_toolchain
    elif java_common.JavaToolchainInfo != platform_common.ToolchainInfo and \
         java_common.JavaToolchainInfo in target:
        toolchain = target[java_common.JavaToolchainInfo]
    else:
        return False
    print("collecting java toolchain info!")

    javac_jars = []
    if hasattr(toolchain, "tools"):
        javac_jars = [
            artifact_location(f)
            for f in toolchain.tools.to_list()
            if f.basename.endswith(".jar")
        ]
    ide_info["java_toolchain_ide_info"] = struct_omit_none(
        javac_jars = javac_jars,
        source_version = toolchain.source_version,
        target_version = toolchain.target_version,
    )

    # update_sync_output_groups(output_groups, "intellij-info-java", depset([ide_info_file]))
    return True

def collect_java_info(ctx, target, feature_configuration, cc_toolchain, ide_info, jar_index_files):
    java = get_java_provider(target)
    if not java:
        return False
    if hasattr(java, "java_outputs") and java.java_outputs:
        java_outputs = java.java_outputs
    elif hasattr(java, "outputs") and java.outputs:
        java_outputs = java.outputs.jars
    else:
        return False

    java_semantics = None

    # java_semantics = semantics.java if hasattr(semantics, "java") else None
    if java_semantics and java_semantics.skip_target(target, ctx):
        return False

    sources = sources_from_target(ctx)
    jars = [library_artifact(output) for output in java_outputs]
    class_jars = [output.class_jar for output in java_outputs if output and output.class_jar]
    output_jars = [jar for output in java_outputs for jar in jars_from_output(output)]
    resolve_files = output_jars
    compile_files = class_jars

    gen_jars = []
    for generated_class_jar, generated_source_jar in _collect_generated_files(java):
        gen_jars.append(annotation_processing_jars(generated_class_jar, generated_source_jar))
        resolve_files += [
            jar
            for jar in [
                generated_class_jar,
                generated_source_jar,
            ]
            if jar != None and not jar.is_source
        ]
        compile_files += [
            jar
            for jar in [generated_class_jar]
            if jar != None and not jar.is_source
        ]

    jdeps = None
    jdeps_file = None
    if java_semantics and hasattr(java_semantics, "get_filtered_jdeps"):
        jdeps_file = java_semantics.get_filtered_jdeps(target)
    if jdeps_file == None and hasattr(java, "outputs") and hasattr(java.outputs, "jdeps") and java.outputs.jdeps:
        jdeps_file = java.outputs.jdeps
    if jdeps_file:
        jdeps = artifact_location(jdeps_file)
        resolve_files.append(jdeps_file)

    java_sources, gen_java_sources, srcjars = divide_java_sources(ctx)

    if java_semantics:
        srcjars = java_semantics.filter_source_jars(target, ctx, srcjars)

    # filtered_gen_jar = None
    # if java_sources and (gen_java_sources or srcjars):
    #     filtered_gen_jar, filtered_gen_resolve_files = _build_filtered_gen_jar(
    #         ctx,
    #         target,
    #         java_outputs,
    #         gen_java_sources,
    #         srcjars,
    #     )
    #     resolve_files += filtered_gen_resolve_files

    # Custom lint checks are incorporated as java plugins. We collect them here and register them with the IDE so that the IDE can also run the same checks.
    plugin_processor_jars = []
    if hasattr(java, "annotation_processing") and java.annotation_processing:
        plugin_processor_jar_files = java.annotation_processing.processor_classpath.to_list()
        resolve_files += plugin_processor_jar_files
        plugin_processor_jars = [annotation_processing_jars(jar, None) for jar in plugin_processor_jar_files]

    if java_outputs:
        class_jars = [info.class_jar for info in java_outputs]
        for jar in class_jars:
            jar_index_file = build_jar_index(ctx, target, jar)
            jar_index_files.append(jar_index_file)

    java_info = struct_omit_none(
        # filtered_gen_jar = filtered_gen_jar,
        generated_jars = gen_jars,
        jars = jars,
        jdeps = jdeps,
        main_class = getattr(ctx.rule.attr, "main_class", None),
        # package_manifest = artifact_location(package_manifest),
        sources = sources,
        test_class = getattr(ctx.rule.attr, "test_class", None),
        # TODO(b/211509545): re-enable plugin_processor_jars after mac user get latest aswb/ ij (2021.12.14.0.1 or later)
        #plugin_processor_jars = plugin_processor_jars,
    )

    # print("java_info:", java_info)
    ide_info["java_ide_info"] = java_info

    # java_info_files.append(ide_info_file)
    # update_sync_output_groups(output_groups, "intellij-info-java", depset(ide_info_files))
    # update_sync_output_groups(output_groups, "intellij-compile-java", depset(compile_files))
    # update_sync_output_groups(output_groups, "intellij-resolve-java", depset(resolve_files))

    # also add transitive hjars + src jars, to catch implicit deps
    # if hasattr(java, "transitive_compile_time_jars"):
    #     update_set_in_dict(output_groups, "intellij-resolve-java-direct-deps", java.transitive_compile_time_jars)
    #     update_set_in_dict(output_groups, "intellij-resolve-java-direct-deps", java.transitive_source_jars)
    return True

def _java_indexer_aspect_impl(target, ctx):
    # print("indexer aspect visiting:", str(target.label))
    deps = []
    if hasattr(ctx.rule.attr, "deps"):
        deps.extend(ctx.rule.attr.deps)
    if hasattr(ctx.rule.attr, "runtime_deps"):
        deps.extend(ctx.rule.attr.runtime_deps)

    transitive_info_file = []
    transitive_jar_index_files = []
    java_info_files = []
    for dep in deps:
        if JarIndexerAspect not in dep:
            continue
        transitive_info_file.append(dep[JarIndexerAspect].info_file)
        transitive_jar_index_files.append(dep[JarIndexerAspect].jar_index_files)
        java_info_files.append(dep[OutputGroupInfo].java_info_files)

    # We support only these rule kinds.
    if ctx.rule.kind not in _all_rules:
        return [
            JarIndexerAspect(
                info_file = depset(transitive = transitive_info_file),
                jar_index_files = depset(transitive = transitive_jar_index_files),
            ),
            OutputGroupInfo(
                java_info_files = depset(transitive = java_info_files),
            ),
        ]

    cc_toolchain = find_cpp_toolchain(ctx)
    feature_configuration = cc_common.configure_features(
        ctx = ctx,
        cc_toolchain = cc_toolchain,
        requested_features = ctx.features,
        unsupported_features = ctx.disabled_features + UNSUPPORTED_FEATURES,
    )

    rule_attrs = ctx.rule.attr

    # Collect direct dependencies
    direct_dep_targets = collect_targets_from_attrs(
        rule_attrs,
        DEPS,
        # semantics_extra_deps(DEPS, semantics, "extra_deps"),
    )
    direct_deps = make_deps(direct_dep_targets, COMPILE_TIME)

    # Add exports from direct dependencies
    exported_deps_from_deps = []
    for dep in direct_dep_targets:
        exported_deps_from_deps = exported_deps_from_deps + dep.intellij_info.export_deps

    # Combine into all compile time deps
    compiletime_deps = direct_deps + exported_deps_from_deps

    # Propagate my own exports
    export_deps = []
    direct_exports = []
    if JavaInfo in target:
        direct_exports = collect_targets_from_attrs(rule_attrs, ["exports"])
        export_deps.extend(make_deps(direct_exports, COMPILE_TIME))

        # Collect transitive exports
        for export in direct_exports:
            export_deps.extend(export.intellij_info.export_deps)

        if ctx.rule.kind == "android_library":
            # Empty android libraries export all their dependencies.
            if not hasattr(rule_attrs, "srcs") or not ctx.rule.attr.srcs:
                export_deps.extend(compiletime_deps)

        # Deduplicate the entries
        export_deps = depset(export_deps).to_list()

    # runtime_deps
    runtime_dep_targets = collect_targets_from_attrs(
        rule_attrs,
        RUNTIME_DEPS,
    )
    runtime_deps = make_deps(runtime_dep_targets, RUNTIME)
    all_deps = depset(compiletime_deps + runtime_deps).to_list()

    # extra prerequisites
    extra_prerequisite_targets = collect_targets_from_attrs(
        rule_attrs,
        PREREQUISITE_DEPS,
        # semantics_extra_deps(PREREQUISITE_DEPS, semantics, "extra_prerequisites"),
    )

    forwarded_deps = _get_forwarded_deps(target, ctx) + direct_exports

    # Roll up output files from my prerequisites
    prerequisites = direct_dep_targets + runtime_dep_targets + extra_prerequisite_targets + direct_exports
    output_groups = dict()
    for dep in prerequisites:
        for k, v in dep.intellij_info.output_groups.items():
            if dep in forwarded_deps:
                # unconditionally roll up deps for these targets
                output_groups[k] = output_groups[k] + [v] if k in output_groups else [v]
                continue

            # roll up outputs of direct deps into '-direct-deps' output group
            if k.endswith("-direct-deps"):
                continue
            if k.endswith("-outputs"):
                directs = k[:-len("outputs")] + "direct-deps"
                output_groups[directs] = output_groups[directs] + [v] if directs in output_groups else [v]
                continue

            # everything else gets rolled up transitively
            output_groups[k] = output_groups[k] + [v] if k in output_groups else [v]

    # Convert output_groups from lists to depsets after the lists are finalized. This avoids
    # creating and growing depsets gradually, as that results in depsets many levels deep:
    # a construct which would give the build system some trouble.
    for k, v in output_groups.items():
        output_groups[k] = depset(transitive = output_groups[k])

    # Initialize the ide info dict, and corresponding output file
    # This will be passed to each language-specific handler to fill in as required
    file_name = target.label.name

    # bazel allows target names differing only by case, so append a hash to support
    # case-insensitive file systems
    file_name = file_name + "-" + str(hash(file_name))
    aspect_ids = get_aspect_ids(ctx, target)
    if aspect_ids:
        aspect_hash = hash(".".join(aspect_ids))
        file_name = file_name + "-" + str(aspect_hash)
    file_name = file_name + ".java_info.json"
    ide_info_file = ctx.actions.declare_file(file_name)

    output_groups = dict()
    target_key = make_target_key(target.label, aspect_ids)
    ide_info = dict(
        # build_file_artifact_location = build_file_artifact_location(ctx),
        features = ctx.features,
        key = target_key,
        kind_string = ctx.rule.kind,
        tags = ctx.rule.attr.tags,
        deps = list(all_deps),
    )

    jar_index_files = []
    if ctx.rule.kind in _java_rules:
        handled = False
        handled = collect_java_info(ctx, target, feature_configuration, cc_toolchain, ide_info, jar_index_files)
        handled = collect_java_toolchain_info(target, ide_info, ide_info_file) or handled

        # elif ctx.rule.kind in _objc_rules:
        #     compile_commands = _objc_compile_commands(ctx, target, feature_configuration, cc_toolchain)

    else:
        fail("unsupported rule: " + ctx.rule.kind)

    # Write the commands for this target.
    info = struct_omit_none(**ide_info)
    ctx.actions.write(
        content = json.encode(info),
        output = ide_info_file,
    )
    ctx.actions.write(ide_info_file, info.to_json())

    info_file = depset([ide_info_file], transitive = transitive_info_file)

    return [
        JarIndexerAspect(
            info_file = info_file,
            jar_index_files = depset(direct = jar_index_files, transitive = transitive_jar_index_files),
        ),
        OutputGroupInfo(
            java_info_files = depset([ide_info_file], transitive = java_info_files),
            jar_index_files = depset(jar_index_files, transitive = transitive_jar_index_files),
        ),
    ]

java_indexer_aspect = aspect(
    attr_aspects = ["deps", "runtime_deps"],
    attrs = {
        "_cc_toolchain": attr.label(
            default = Label("@bazel_tools//tools/cpp:current_cc_toolchain"),
        ),
        # "_jarindexer": attr.label(
        #     default = Label("//cmd/jarindexer"),
        #     cfg = "exec",
        #     executable = True,
        # ),
        "_jarindexer2": attr.label(
            default = Label("//cmd/jarindexer:jarindexer2"),
            cfg = "exec",
            executable = True,
        ),
        "_ijar": attr.label(
            default = Label("@bazel_tools//tools/jdk:ijar"),
            executable = True,
            cfg = "exec",
            allow_files = True,
        ),
    },
    fragments = ["cpp", "java"],
    provides = [JarIndexerAspect],
    toolchains = ["@bazel_tools//tools/cpp:toolchain_type"],
    implementation = _java_indexer_aspect_impl,
    apply_to_generating_rules = True,
)
