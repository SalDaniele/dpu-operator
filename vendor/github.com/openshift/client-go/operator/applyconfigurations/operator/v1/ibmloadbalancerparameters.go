// Code generated by applyconfiguration-gen. DO NOT EDIT.

package v1

import (
	v1 "github.com/openshift/api/operator/v1"
)

// IBMLoadBalancerParametersApplyConfiguration represents an declarative configuration of the IBMLoadBalancerParameters type for use
// with apply.
type IBMLoadBalancerParametersApplyConfiguration struct {
	Protocol *v1.IngressControllerProtocol `json:"protocol,omitempty"`
}

// IBMLoadBalancerParametersApplyConfiguration constructs an declarative configuration of the IBMLoadBalancerParameters type for use with
// apply.
func IBMLoadBalancerParameters() *IBMLoadBalancerParametersApplyConfiguration {
	return &IBMLoadBalancerParametersApplyConfiguration{}
}

// WithProtocol sets the Protocol field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the Protocol field is set to the value of the last call.
func (b *IBMLoadBalancerParametersApplyConfiguration) WithProtocol(value v1.IngressControllerProtocol) *IBMLoadBalancerParametersApplyConfiguration {
	b.Protocol = &value
	return b
}
