load("@io_bazel_rules_go//go:def.bzl", "go_binary")

go_binary(
    name = "client",
    srcs = ["main.go", "newfile.go", "math.go"],
    deps = ["//demo:mrpc"],
    visibility = ["//visibility:public"],
)