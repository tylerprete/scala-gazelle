load("@bazel_tools//tools/build_defs/repo:http.bzl", "http_archive", "http_file")

def _maybe(repo_rule, name, **kwargs):
    if name not in native.existing_rules():
        repo_rule(name = name, **kwargs)

def workspace_deps():
    io_bazel_rules_go()  # via bazel_gazelle

    bazel_gazelle()  # via <TOP>
    rules_proto()  # via <TOP>
    build_stack_rules_proto()
    build_bazel_rules_nodejs()  # via <TOP>
    rules_jvm_external()
    io_bazel_rules_scala()
    protobuf_core_deps()
    viz_js_lite()

def io_bazel_rules_go():
    # Release: v0.35.0
    # TargetCommitish: release-0.35
    # Date: 2022-09-11 15:59:49 +0000 UTC
    # URL: https://github.com/bazelbuild/rules_go/releases/tag/v0.35.0
    # Size: 931734 (932 kB)
    _maybe(
        http_archive,
        name = "io_bazel_rules_go",
        sha256 = "cc027f11f98aef8bc52c472ced0714994507a16ccd3a0820b2df2d6db695facd",
        strip_prefix = "rules_go-0.35.0",
        urls = ["https://github.com/bazelbuild/rules_go/archive/v0.35.0.tar.gz"],
    )

def bazel_gazelle():
    # Branch: master
    # Commit: c4ec7765fad672d87548bdc43e740fb5843f0839
    # Date: 2022-10-18 17:28:02 +0000 UTC
    # URL: https://github.com/bazelbuild/bazel-gazelle/commit/c4ec7765fad672d87548bdc43e740fb5843f0839
    #
    # Add size argument to `gazelle_generation_test` (#1351)
    #
    # * Add size argument to gazelle integration test macro
    # Size: 1573804 (1.6 MB)
    _maybe(
        http_archive,
        name = "bazel_gazelle",
        patch_args = ["-p1"],
        patches = ["//third_party/bazelbuild/bazel-gazelle:rule-attrassignment-api.patch"],
        sha256 = "33ad1ec6020e6660e921cab57ea8d49009af1b26b2434c986930ddab620feeb7",
        strip_prefix = "bazel-gazelle-c4ec7765fad672d87548bdc43e740fb5843f0839",
        urls = ["https://github.com/bazelbuild/bazel-gazelle/archive/c4ec7765fad672d87548bdc43e740fb5843f0839.tar.gz"],
    )

def local_bazel_gazelle():
    _maybe(
        native.local_repository,
        name = "bazel_gazelle",
        path = "/Users/i868039/go/src/github.com/bazelbuild/bazel-gazelle",
    )

def rules_proto():
    _maybe(
        http_archive,
        name = "rules_proto",
        sha256 = "9fc210a34f0f9e7cc31598d109b5d069ef44911a82f507d5a88716db171615a8",
        strip_prefix = "rules_proto-f7a30f6f80006b591fa7c437fe5a951eb10bcbcf",
        urls = [
            "https://github.com/bazelbuild/rules_proto/archive/f7a30f6f80006b591fa7c437fe5a951eb10bcbcf.tar.gz",
        ],
    )

def build_stack_rules_proto():
    # Release: v2.0.1
    # TargetCommitish: master
    # Date: 2022-10-20 02:38:27 +0000 UTC
    # URL: https://github.com/stackb/rules_proto/releases/tag/v2.0.1
    # Size: 2071295 (2.1 MB)
    http_archive(
        name = "build_stack_rules_proto",
        sha256 = "ac7e2966a78660e83e1ba84a06db6eda9a7659a841b6a7fd93028cd8757afbfb",
        strip_prefix = "rules_proto-2.0.1",
        urls = ["https://github.com/stackb/rules_proto/archive/v2.0.1.tar.gz"],
    )

def build_bazel_rules_nodejs():
    _maybe(
        http_archive,
        name = "build_bazel_rules_nodejs",
        sha256 = "4501158976b9da216295ac65d872b1be51e3eeb805273e68c516d2eb36ae1fbb",
        urls = [
            "https://github.com/bazelbuild/rules_nodejs/releases/download/4.4.1/rules_nodejs-4.4.1.tar.gz",
        ],
    )

def rules_jvm_external():
    _maybe(
        http_archive,
        name = "rules_jvm_external",
        sha256 = "31701ad93dbfe544d597dbe62c9a1fdd76d81d8a9150c2bf1ecf928ecdf97169",
        strip_prefix = "rules_jvm_external-4.0",
        urls = [
            "https://github.com/bazelbuild/rules_jvm_external/archive/4.0.zip",
        ],
    )

def io_bazel_rules_scala():
    _maybe(
        http_archive,
        name = "io_bazel_rules_scala",
        sha256 = "0701ee4e1cfd59702d780acde907ac657752fbb5c7d08a0ec6f58ebea8cd0efb",
        strip_prefix = "rules_scala-2437e40131072cadc1628726775ff00fa3941a4a",
        urls = [
            "https://github.com/bazelbuild/rules_scala/archive/2437e40131072cadc1628726775ff00fa3941a4a.tar.gz",
        ],
    )

def protobuf_core_deps():
    bazel_skylib()  # via com_google_protobuf
    rules_python()  # via com_google_protobuf
    zlib()  # via com_google_protobuf
    com_google_protobuf()  # via <TOP>

def bazel_skylib():
    _maybe(
        http_archive,
        name = "bazel_skylib",
        sha256 = "ebdf850bfef28d923a2cc67ddca86355a449b5e4f38b0a70e584dc24e5984aa6",
        strip_prefix = "bazel-skylib-f80bc733d4b9f83d427ce3442be2e07427b2cc8d",
        urls = [
            "https://github.com/bazelbuild/bazel-skylib/archive/f80bc733d4b9f83d427ce3442be2e07427b2cc8d.tar.gz",
        ],
    )

def rules_python():
    _maybe(
        http_archive,
        name = "rules_python",
        sha256 = "8cc0ad31c8fc699a49ad31628273529ef8929ded0a0859a3d841ce711a9a90d5",
        strip_prefix = "rules_python-c7e068d38e2fec1d899e1c150e372f205c220e27",
        urls = [
            "https://github.com/bazelbuild/rules_python/archive/c7e068d38e2fec1d899e1c150e372f205c220e27.tar.gz",
        ],
    )

def zlib():
    _maybe(
        http_archive,
        name = "zlib",
        sha256 = "c3e5e9fdd5004dcb542feda5ee4f0ff0744628baf8ed2dd5d66f8ca1197cb1a1",
        strip_prefix = "zlib-1.2.11",
        urls = [
            "https://mirror.bazel.build/zlib.net/zlib-1.2.11.tar.gz",
            "https://zlib.net/zlib-1.2.11.tar.gz",
        ],
        build_file = "@build_stack_rules_proto//third_party:zlib.BUILD",
    )

def com_google_protobuf():
    _maybe(
        http_archive,
        name = "com_google_protobuf",
        sha256 = "d0f5f605d0d656007ce6c8b5a82df3037e1d8fe8b121ed42e536f569dec16113",
        strip_prefix = "protobuf-3.14.0",
        urls = [
            "https://github.com/protocolbuffers/protobuf/archive/v3.14.0.tar.gz",
        ],
    )

def viz_js_lite():
    # HTTP/2.0 200 OK
    # Content-Type: application/javascript; charset=utf-8
    # Date: Sat, 26 Mar 2022 00:35:46 GMT
    # Last-Modified: Mon, 04 May 2020 16:17:44 GMT
    # Size: 1439383 (1.4 MB)
    _maybe(
        http_file,
        name = "cdnjs_cloudflare_com_ajax_libs_viz_js_2_1_2_lite_render_js",
        sha256 = "1344fd45812f33abcb3de9857ebfdd599e57f49e3d0849841e75e28be1dd6959",
        urls = ["https://cdnjs.cloudflare.com/ajax/libs/viz.js/2.1.2/lite.render.js"],
    )
