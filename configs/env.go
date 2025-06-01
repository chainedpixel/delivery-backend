package config

import (
	"errors"
	"fmt"
	errPackage "github.com/MarlonG1/delivery-backend/configs/error"
	"github.com/spf13/viper"
	"path/filepath"
	"runtime"
	"strings"
)

type EnvConfig struct {
	Server struct {
		Port      string
		JWTSecret string
		Debug     bool
	}
	Database struct {
		Host     string
		Port     string
		Name     string
		User     string
		Password string
		Charset  string
		Driver   string
	}
	Redis struct {
		Host     string
		Port     string
		Password string
	}
	Log struct {
		Level       string
		FileLogging bool
	}
	EmailConfig struct {
		SMTPHost      string `json:"smtp_host"`
		SMTPPort      int    `json:"smtp_port"`
		SMTPUsername  string `json:"smtp_username"`
		SMTPPassword  string `json:"smtp_password"`
		FromEmail     string `json:"from_email"`
		FromName      string `json:"from_name"`
		EnableTLS     bool   `json:"enable_tls"`
		SkipTLSVerify bool   `json:"skip_tls_verify"`
		Timeout       int    `json:"timeout"`
		RetryAttempts int    `json:"retry_attempts"`
		EnableEmails  bool   `json:"enable_emails"`
	}
}

func NewEnvConfig() (*EnvConfig, error) {
	v := viper.New()

	_, b, _, _ := runtime.Caller(0)
	projectRoot := filepath.Join(filepath.Dir(b), "..")

	v.SetConfigName(".env")
	v.SetConfigType("env")
	v.AddConfigPath(projectRoot)
	fmt.Println("📁 Project Root:", projectRoot)

	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if err := v.ReadInConfig(); err != nil {
		var configFileNotFoundError viper.ConfigFileNotFoundError
		if !errors.As(err, &configFileNotFoundError) {
			return nil, errPackage.ErrEnvFileNotFound
		}
		return nil, err
	}

	v.AutomaticEnv()
	MapEnvKeys(v)

	var config EnvConfig
	if err := v.Unmarshal(&config); err != nil {
		return nil, errPackage.ErrFailedToLoadEnv
	}

	// Debug de configuración de email
	debugEmailConfig(&config)

	return &config, validateConfig(&config)
}

func validateConfig(config *EnvConfig) error {
	// Validaciones de base de datos
	if config.Database.Host == "" {
		return fmt.Errorf("DB_HOST is required")
	}
	if config.Database.Password == "" {
		return fmt.Errorf("DB_PASSWORD is required")
	}
	if config.Database.User == "" {
		return fmt.Errorf("DB_USER is required")
	}

	// Validaciones de email (si está habilitado)
	if config.EmailConfig.EnableEmails {
		if config.EmailConfig.SMTPHost == "" {
			return fmt.Errorf("SMTP_HOST is required when emails are enabled")
		}
		if config.EmailConfig.SMTPPort == 0 {
			return fmt.Errorf("SMTP_PORT is required and cannot be 0")
		}
		if config.EmailConfig.SMTPUsername == "" {
			return fmt.Errorf("SMTP_USERNAME is required when emails are enabled")
		}
		if config.EmailConfig.SMTPPassword == "" {
			return fmt.Errorf("SMTP_PASSWORD is required when emails are enabled")
		}
		if config.EmailConfig.FromEmail == "" {
			return fmt.Errorf("FROM_EMAIL is required when emails are enabled")
		}
	}

	return nil
}

func MapEnvKeys(v *viper.Viper) {
	// .env keys for database connection
	v.Set("database.host", v.GetString("db_host"))
	v.Set("database.port", v.GetString("db_port"))
	v.Set("database.name", v.GetString("db_name"))
	v.Set("database.user", v.GetString("db_user"))
	v.Set("database.password", v.GetString("db_password"))
	v.Set("database.charset", v.GetString("db_charset"))
	v.Set("database.driver", v.GetString("db_driver"))

	// .env keys for redis connection
	v.Set("redis.host", v.GetString("redis_host"))
	v.Set("redis.port", v.GetString("redis_port"))
	v.Set("redis.password", v.GetString("redis_password"))

	// .env keys for server configuration
	v.Set("server.port", v.GetString("server_port"))
	v.Set("server.jwtSecret", v.GetString("jwt_secret"))
	v.Set("server.debug", v.GetBool("debug"))

	// .env keys for log configuration
	v.Set("log.level", v.GetString("log_level"))
	v.Set("log.fileLogging", v.GetString("log_file_logging"))

	// ✅ AGREGAR: .env keys for email configuration (¡ESTO FALTABA!)
	v.Set("emailconfig.smtphost", v.GetString("smtp_host"))
	v.Set("emailconfig.smtpport", v.GetInt("smtp_port"))
	v.Set("emailconfig.smtpusername", v.GetString("smtp_username"))
	v.Set("emailconfig.smtppassword", v.GetString("smtp_password"))
	v.Set("emailconfig.fromemail", v.GetString("from_email"))
	v.Set("emailconfig.fromname", v.GetString("from_name"))
	v.Set("emailconfig.enabletls", v.GetBool("smtp_enable_tls"))
	v.Set("emailconfig.skiptlsverify", v.GetBool("smtp_skip_tls_verify"))
	v.Set("emailconfig.timeout", v.GetInt("email_timeout"))
	v.Set("emailconfig.retryattempts", v.GetInt("email_retry_attempts"))
	v.Set("emailconfig.enableemails", v.GetBool("enable_emails"))
}

// debugEmailConfig muestra la configuración de email para debugging
func debugEmailConfig(config *EnvConfig) {
	println("📧 ===== EMAIL CONFIGURATION DEBUG =====")
	println("SMTP Host:", config.EmailConfig.SMTPHost)
	println("SMTP Port:", config.EmailConfig.SMTPPort)
	println("SMTP Username:", config.EmailConfig.SMTPUsername)
	println("SMTP Password:", maskPassword(config.EmailConfig.SMTPPassword))
	println("From Email:", config.EmailConfig.FromEmail)
	println("From Name:", config.EmailConfig.FromName)
	println("Enable TLS:", config.EmailConfig.EnableTLS)
	println("Skip TLS Verify:", config.EmailConfig.SkipTLSVerify)
	println("Enable Emails:", config.EmailConfig.EnableEmails)
	println("Timeout:", config.EmailConfig.Timeout)
	println("Retry Attempts:", config.EmailConfig.RetryAttempts)
	println("========================================")
}

// maskPassword enmascara la contraseña para el debug
func maskPassword(password string) string {
	if len(password) <= 4 {
		return "****"
	}
	return password[:2] + "****" + password[len(password)-2:]
}
