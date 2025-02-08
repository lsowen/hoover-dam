package cmd

import (
	"context"
	"fmt"
	"time"

	"github.com/lsowen/hoover-dam/pkg/db"
	"github.com/spf13/cobra"
	"github.com/treeverse/lakefs/pkg/auth/keys"
)

var superuserCmd = &cobra.Command{
	Use:   "superuser",
	Short: "Create users with admin credentials",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := loadConfig()
		if err != nil {
			return fmt.Errorf("loading superuser command config: %w", err)
		}

		username, err := cmd.Flags().GetString("user-name")
		if err != nil {
			return fmt.Errorf("getting --user-name: %w", err)
		}

		database, err := db.NewDatabase(cmd.Context(), *cfg)
		if err != nil {
			return fmt.Errorf("initializing database: %w", err)
		}
		_, credential, err := CreateAdminUser(cmd.Context(), *database, username)
		if err != nil {
			return fmt.Errorf("creating admin user: %w", err)
		}

		fmt.Printf("credentials:\n  access_key_id: %s\n  secret_access_key: %s\n",
			credential.AccessKeyId, credential.SecretAccessKey)

		return nil
	},
}

func CreateAdminUser(ctx context.Context, database db.Database, username string) (*db.User, *db.Credential, error) {

	user, err := database.GetUser(ctx, username)
	if err != nil {
		return nil, nil, fmt.Errorf("getting user from the db: %w", err)
	}

	if user == nil {
		user = &db.User{
			Username:     username,
			CreationDate: time.Now(),
		}
		err = database.CreateUser(ctx, user)
		if err != nil {
			return nil, nil, fmt.Errorf("creating new user: %w", err)
		}
	}

	err = database.AddGroupMember(ctx, "Admins", username)
	if err != nil {
		return nil, nil, fmt.Errorf("adding user to admins group: %w", err)
	}

	accessKeyID := keys.GenAccessKeyID()
	secretAccessKey := keys.GenSecretAccessKey()

	credential, err := database.CreateUserCredential(ctx, username, accessKeyID, secretAccessKey)
	if err != nil {
		return nil, nil, fmt.Errorf("creating user credential: %w", err)
	}
	return user, &credential, nil
}

func init() {
	rootCmd.AddCommand(superuserCmd)
	flags := superuserCmd.Flags()
	flags.String("user-name", "", "identifier for the user")
	superuserCmd.MarkFlagRequired("user-name")
}
