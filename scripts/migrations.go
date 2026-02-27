package scripts

import (
	"log"
	"os"
	"os/exec"
)

func RunMigrations(dsn string) {
	cmd := exec.Command("make", "migrate-dev")
	cmd.Env = append(os.Environ(), "DATABASE_URL="+dsn) // pass DSN dynamically

	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatalf("❌ Migration failed: %v\n%s", err, string(output))
	}
	log.Println("✅ Database migrated successfully")
}
