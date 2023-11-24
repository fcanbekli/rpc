load("@io_bazel_rules_go//go:def.bzl", "go_binary", "go_library")

go_binary(
    name = "server",
    srcs = ["main.go", "newfile.go", "math.go"],
    deps = ["//demo:mrpc"],
    visibility = ["//visibility:public"],
)