package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/gov4git/gov4git/v2/proto/docket/ops"
	"github.com/gov4git/gov4git/v2/proto/docket/policy"
	"github.com/gov4git/gov4git/v2/proto/docket/schema"
	"github.com/gov4git/gov4git/v2/proto/member"
	"github.com/gov4git/lib4git/form"
	"github.com/spf13/cobra"
)

var (
	motionCmd = &cobra.Command{
		Use:   "motion",
		Short: "Manage motions (concerns and proposals)",
		Long:  ``,
		Run:   func(cmd *cobra.Command, args []string) {},
	}

	motionOpenCmd = &cobra.Command{
		Use:   "open",
		Short: "Open a new motion",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
			LoadConfig()
			ops.OpenMotion(
				ctx,
				setup.Organizer,
				schema.MotionID(motionName),
				schema.ParseMotionType(ctx, motionType),
				schema.PolicyName(motionPolicy),
				member.User(motionAuthor),
				motionTitle,
				motionDesc,
				motionTrackerURL,
				nil,
			)
			// fmt.Fprint(os.Stdout, form.SprintJSON(chg.Result))
		},
	}

	motionCloseCmd = &cobra.Command{
		Use:   "close",
		Short: "Close a motion",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
			LoadConfig()
			var d schema.Decision
			if motionAccept {
				d = schema.Accept
			} else {
				d = schema.Reject
			}
			ops.CloseMotion(
				ctx,
				setup.Organizer,
				schema.MotionID(motionName),
				d,
			)
			// fmt.Fprint(os.Stdout, form.SprintJSON(chg.Result))
		},
	}

	motionListCmd = &cobra.Command{
		Use:   "list",
		Short: "List motions",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
			LoadConfig()
			l := ops.ListMotionViews(ctx, setup.Gov)
			fmt.Fprint(os.Stdout, form.SprintJSON(l))
		},
	}

	motionShowCmd = &cobra.Command{
		Use:   "show",
		Short: "Show motion state and associated objects",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
			LoadConfig()
			s := ops.ShowMotion(ctx, setup.Gov, schema.MotionID(motionName))
			fmt.Fprint(os.Stdout, form.SprintJSON(s))
		},
	}
)

var (
	motionName       string
	motionPolicy     string
	motionAuthor     string
	motionTitle      string
	motionDesc       string
	motionType       string
	motionTrackerURL string
	motionAccept     bool
)

func init() {
	motionCmd.AddCommand(motionOpenCmd)
	motionOpenCmd.Flags().StringVar(&motionName, "name", "", "unique name for motion")
	motionOpenCmd.MarkFlagRequired("name")
	motionOpenCmd.Flags().StringVar(&motionPolicy, "policy", "", "policy ("+strings.Join(policy.InstalledMotionPolicies(), ", ")+")")
	motionOpenCmd.MarkFlagRequired("policy")
	motionOpenCmd.Flags().StringVar(&motionAuthor, "author", "", "author user name")
	motionOpenCmd.MarkFlagRequired("author")

	motionOpenCmd.Flags().StringVar(&motionTitle, "title", "", "title for motion")
	motionOpenCmd.Flags().StringVar(&motionDesc, "desc", "", "description for motion")
	motionOpenCmd.Flags().StringVar(&motionType, "type", "concern", "type of motion (concern, proposal)")
	motionOpenCmd.Flags().StringVar(&motionTrackerURL, "tracking", "", "tracking URL for motion")

	motionCmd.AddCommand(motionCloseCmd)
	motionCloseCmd.Flags().StringVar(&motionName, "name", "", "name of motion")
	motionCloseCmd.MarkFlagRequired("name")
	motionCloseCmd.Flags().BoolVar(&motionAccept, "accept", false, "accept/reject")

	motionCmd.AddCommand(motionListCmd)

	motionCmd.AddCommand(motionShowCmd)
	motionShowCmd.Flags().StringVar(&motionName, "name", "", "name of motion")
	motionShowCmd.MarkFlagRequired("name")
}
