package raw

import (
	"fmt"
	"strings"

	"github.com/ctr2cloud/ctr2cloud/pkg/generic/compute"
	"github.com/ctr2cloud/ctr2cloud/pkg/providers/auto"
	"github.com/spf13/cobra"
)

const (
	providerFlag = "provider"
	instanceFlag = "instance"
	noStreamFlag = "no-stream"
)

func init() {
	pFlags := Cmd.PersistentFlags()
	pFlags.StringP(providerFlag, "p", "", "provider to use")
	Cmd.MarkPersistentFlagRequired(providerFlag)

	execFlags := execCmd.Flags()
	execFlags.StringP(instanceFlag, "i", "", "instance to execute command on")
	execCmd.MarkFlagRequired(instanceFlag)
	execFlags.Bool(noStreamFlag, false, "execute command without streaming output")

	Cmd.AddCommand(listCmd)
	Cmd.AddCommand(createCmd)
	Cmd.AddCommand(deleteCmd)
	Cmd.AddCommand(execCmd)
}

var Cmd = &cobra.Command{
	Use:   "raw",
	Short: "raw provides the ability to directly interact with providers",
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "list instances",
	RunE: func(cmd *cobra.Command, args []string) error {
		providerName, _ := cmd.Flags().GetString(providerFlag)
		provider, ok := auto.Providers[providerName]
		if !ok {
			return fmt.Errorf("provider %q not found", providerName)
		}
		instances, err := provider.List()
		if err != nil {
			return fmt.Errorf("listing instances: %w", err)
		}
		for _, instance := range instances {
			fmt.Printf("Name: %s || Internal Name: %s\n", instance.Id, instance.Name)
		}
		return nil
	},
}

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "create an instance",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		providerName, _ := cmd.Flags().GetString(providerFlag)
		provider, ok := auto.Providers[providerName]
		if !ok {
			return fmt.Errorf("provider %q not found", providerName)
		}
		err := provider.Create(compute.InstanceSpec{
			Name: args[0],
		})
		if err != nil {
			return fmt.Errorf("creating instance: %w", err)
		}
		return nil
	},
}

var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "delete an instance",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		providerName, _ := cmd.Flags().GetString(providerFlag)
		provider, ok := auto.Providers[providerName]
		if !ok {
			return fmt.Errorf("provider %q not found", providerName)
		}
		err := provider.Delete(args[0])
		if err != nil {
			return fmt.Errorf("deleting instance: %w", err)
		}
		return nil
	},
}

var execCmd = &cobra.Command{
	Use:   "exec",
	Short: "execute a command on an instance",
	RunE: func(cmd *cobra.Command, args []string) error {
		providerName, _ := cmd.Flags().GetString(providerFlag)
		provider, ok := auto.Providers[providerName]
		if !ok {
			return fmt.Errorf("provider %q not found", providerName)
		}

		instanceName, _ := cmd.Flags().GetString(instanceFlag)
		noStream, _ := cmd.Flags().GetBool(noStreamFlag)

		executor, err := provider.GetCommandExecutor(instanceName)
		if err != nil {
			return fmt.Errorf("getting command executor: %w", err)
		}
		if noStream {
			res, err := executor.ExecString(cmd.Context(), strings.Join(args, " "))
			if err != nil {
				return fmt.Errorf("exec error: %w", err)
			}
			fmt.Print(res)
			return nil
		}

		resChan := executor.ExecStream(cmd.Context(), strings.Join(args, " "))
		for res := range resChan {
			if res.Error != nil {
				return fmt.Errorf("command error: %w", res.Error)
			}
			fmt.Printf("%s", res.Data)
		}
		return nil
	},
}
