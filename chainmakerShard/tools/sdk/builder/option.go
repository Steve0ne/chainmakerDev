/*
 Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
 SPDX-License-Identifier: Apache-2.0
*/
package builder

type CrossBuildOption func(*crossBuildOptions)

type crossBuildOptions struct {
	ProofKey     string
	ParamOptions []ParamsBuildOption
}

func NewCrossBuildOptions(opts ...CrossBuildOption) *crossBuildOptions {
	options := &crossBuildOptions{}
	for _, o := range opts {
		o(options)
	}
	return options
}

func WithUseProofKey(use bool) CrossBuildOption {
	return func(opts *crossBuildOptions) {
		if use {
			opts.ParamOptions = append(opts.ParamOptions, paramsBuildWithProofKey(use, opts.ProofKey))
		}
	}
}

type paramsBuildOptions struct {
	UseProofKey bool
	ProofKey    string
}

type ParamsBuildOption func(*paramsBuildOptions)

func paramsBuildWithProofKey(use bool, key string) ParamsBuildOption {
	return func(opt *paramsBuildOptions) {
		opt.UseProofKey = use
		opt.ProofKey = key
	}
}

func NewParamsBuildOptions(opts ...ParamsBuildOption) *paramsBuildOptions {
	options := &paramsBuildOptions{}
	for _, o := range opts {
		o(options)
	}
	return options
}
