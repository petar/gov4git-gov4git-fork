package policy

import (
	"context"

	"github.com/gov4git/gov4git/proto/docket/schema"
	"github.com/gov4git/gov4git/proto/gov"
	"github.com/gov4git/gov4git/proto/mod"
	"github.com/gov4git/lib4git/form"
	"github.com/gov4git/lib4git/ns"
)

type Policy interface {
	Name() string

	Open(
		ctx context.Context,
		cloned gov.OwnerCloned,
		motion schema.Motion,
		instancePolicyNS ns.NS,
	)

	Score(
		ctx context.Context,
		cloned gov.OwnerCloned,
		motion schema.Motion,
		instancePolicyNS ns.NS,
	) schema.Score

	Close(
		ctx context.Context,
		cloned gov.OwnerCloned,
		motion schema.Motion,
		instancePolicyNS ns.NS,
	)

	Cancel(
		ctx context.Context,
		cloned gov.OwnerCloned,
		motion schema.Motion,
		instancePolicyNS ns.NS,
	)

	Show(
		ctx context.Context,
		cloned gov.OwnerCloned,
		motion schema.Motion,
		instancePolicyNS ns.NS,
	) form.Map

	AddRefTo(
		ctx context.Context,
		cloned gov.OwnerCloned,
		refType schema.RefType,
		from schema.Motion,
		to schema.Motion,
		fromPolicyNS ns.NS,
		toPolicyNS ns.NS,
	)

	AddRefFrom(
		ctx context.Context,
		cloned gov.OwnerCloned,
		refType schema.RefType,
		from schema.Motion,
		to schema.Motion,
		fromPolicyNS ns.NS,
		toPolicyNS ns.NS,
	)

	RemoveRefTo(
		ctx context.Context,
		cloned gov.OwnerCloned,
		refType schema.RefType,
		from schema.Motion,
		to schema.Motion,
		fromPolicyNS ns.NS,
		toPolicyNS ns.NS,
	)

	RemoveRefFrom(
		ctx context.Context,
		cloned gov.OwnerCloned,
		refType schema.RefType,
		from schema.Motion,
		to schema.Motion,
		fromPolicyNS ns.NS,
		toPolicyNS ns.NS,
	)
}

var policyRegistry = mod.NewModuleRegistry[Policy]()

func Install(ctx context.Context, policy Policy) {
	policyRegistry.Set(ctx, policy.Name(), policy)
}

func Get(ctx context.Context, key string) Policy {
	return policyRegistry.Get(ctx, key)
}

func InstalledPolicies() []string {
	return policyRegistry.Keys()
}

func GetMotionPolicy(ctx context.Context, m schema.Motion) Policy {
	return Get(ctx, m.Policy.String())
}

// MotionPolicyNS returns the private policy namespace for a given motion instance.
func MotionPolicyNS(id schema.MotionID) ns.NS {
	return schema.MotionKV.KeyNS(schema.MotionNS, id).Append("policy")
}
