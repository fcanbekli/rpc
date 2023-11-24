load("@io_bazel_rules_go//go:def.bzl", "go_binary", "go_test", "go_library")

go_binary(
    name = "trpc",
    srcs = [
        "trpc.go",
        "strings.go",
        "size.go",
        "func_desc.go"
    ],
    visibility = ["//visibility:public"],
)

go_library(
    name = "trpc_lib",
    srcs = ["func_desc.go"],
        importpath = "trpc_lib",
    visibility = ["//visibility:public"],
)

go_test(
    name = "trpc_test"
    ,
    srcs = ["func_desc.go", "func_desc_test.go"],
   deps = [":trpc_lib"],
)