package cmd

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/lsowen/hoover-dam/pkg/db"
	"github.com/spf13/cobra"
	"github.com/treeverse/lakefs/pkg/auth/keys"
)

var superuserCmd = &cobra.Command{
	Use:   "superuser",
	Short: "Create users with admin credentials",
	Run: func(cmd *cobra.Command, args []string) {
		cfg := loadConfig()

		username, err := cmd.Flags().GetString("user-name")
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		database, err := db.NewDatabase(cmd.Context(), *cfg)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)

		}
		_, credential, err := CreateAdminUser(cmd.Context(), *database, username)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		fmt.Printf("credentials:\n  access_key_id: %s\n  secret_access_key: %s\n",
			credential.AccessKeyId, credential.SecretAccessKey)
	},
}

func CreateAdminUser(ctx context.Context, database db.Database, username string) (*db.User, *db.Credential, error) {

	user, err := database.GetUser(ctx, username)
	if err != nil {
		return nil, nil, err
	}

	if user == nil {
		user = &db.User{
			Username:     username,
			CreationDate: time.Now(),
		}
		err = database.CreateUser(ctx, user)
		if err != nil {
			return nil, nil, err
		}
	}

	err = database.AddGroupMember(ctx, "Admins", username)
	if err != nil {
		return nil, nil, err
	}

	accessKeyID := keys.GenAccessKeyID()
	secretAccessKey := keys.GenSecretAccessKey()

	credential, err := database.CreateUserCredential(ctx, username, accessKeyID, secretAccessKey)
	if err != nil {
		return nil, nil, err
	}
	return user, &credential, nil
}

func init() {
	rootCmd.AddCommand(superuserCmd)
	flags := superuserCmd.Flags()
	flags.String("user-name", "", "identifier for the user")
	superuserCmd.MarkFlagRequired("user-name")
}
