package main

import (
	"flag"
	"os"
	"slices"
	"strings"
	"text/template"
)

var versionEnv = ServiceEnvironmentValues{
	Key:   "DANK_VERSION",
	Value: "${DANK_VERSION:-${LATEST_TAG:-${COMMIT_SHA}}}",
}

const (
	mariadbImage = "mariadb:11.7"
	redisImage   = "redis:7.2.4"
)

const dockerComposeFullTemplate = `services:
{{- range .Services }}
{{ . }}
{{- end }}

{{- if .Networks }}
networks:
{{- range .Networks }}
{{ . }}
{{- end }}
{{- end }}

{{- if .Volumes }}
volumes:
{{- range .Volumes }}
{{ . }}
{{- end }}
{{- end }}`

type TemplateValues struct {
	ServerPort   string
	WebPort      string
	CdnPort      string
	EventHubPort string
	YtDlPort     string

	NetworkName string

	FilesPath    string
	DbConfigPath string
	DbDataPath   string

	MariaDbImage string
	RedisImage   string

	GoEnv           string
	ExternalNetwork bool
	EnabledServices []string
}

const serviceTemplate = `  {{ .ServiceName }}:
    container_name: "{{ .ServiceName }}"
    {{- if .BuildContext }}
    build:
      context: {{ .BuildContext }}
      dockerfile: {{ .Dockerfile }}
    {{- end }}
    {{- if .ContainerImage }}
    image: "{{ .ContainerImage }}"
    {{- end }}
    restart: "always"
    ports:
    {{- range .Ports }}
      - "{{ .OutPort }}:{{ .ContainerPort }}"
    {{- end }}
    stdin_open: true
    {{- if .Environment}}
    environment:
      {{- range .Environment }}
      {{ .Key }}: "{{ .Value }}"
      {{- end }}
    {{- end }}
    env_file:
      - {{ .EnvFile }}
    {{- if .Volumes}}
    volumes:
      {{- range .Volumes }}
      - {{ .VolumeName }}:{{ .MountPath }}
      {{- end }}
    {{- end }}
    networks:
    {{- range .Networks }}
      - {{ . }}
    {{- end }}
    {{- if .DependsOn }}
    depends_on:
      {{- range .DependsOn }}
      - {{ . }}
      {{- end }}
    {{- end }}
	{{- if .Command }}
    command: >
      {{ .Command }}
	{{- end }}
`

type ServicePortsValues struct {
	OutPort       string
	ContainerPort string
}

type ServiceEnvironmentValues struct {
	Key   string
	Value string
}

type ServiceVolumesValues struct {
	VolumeName string
	MountPath  string
}

type ServiceValues struct {
	ServiceName    string
	BuildContext   string
	Dockerfile     string
	ContainerImage string
	Ports          []ServicePortsValues
	Environment    []ServiceEnvironmentValues
	EnvFile        string
	Volumes        []ServiceVolumesValues
	Networks       []string
	DependsOn      []string
	Command        string
}

func getService(values ServiceValues) (string, error) {
	out := new(strings.Builder)

	t := template.Must(template.New(values.ServiceName).Parse(serviceTemplate))
	err := t.Execute(out, values)
	if err != nil {
		return "", err
	}

	return out.String(), nil
}

const volumeTemplate = `  {{ .VolumeName }}:
    driver: local
    driver_opts:
      type: none
      o: bind
      device: {{ .LocalPath }}
`

type VolumeValues struct {
	VolumeName string
	LocalPath  string
}

func getVolume(values VolumeValues) (string, error) {
	out := new(strings.Builder)

	t := template.Must(template.New(values.VolumeName).Parse(volumeTemplate))
	err := t.Execute(out, values)
	if err != nil {
		return "", err
	}

	return out.String(), nil
}

var networksTemplates = map[string]string{
	"internal-network": `  {{ .NetworkName }}: {}
`,
	"external-network": `  {{ .NetworkName }}:
    external: true
`,
}

func getInternalNetwork(name string) (string, error) {
	out := new(strings.Builder)

	t := template.Must(template.New(name).Parse(networksTemplates["internal-network"]))
	err := t.Execute(out, map[string]string{
		"NetworkName": name,
	})
	if err != nil {
		return "", err
	}

	return out.String(), nil
}

func getExternalNetwork(name string) (string, error) {
	out := new(strings.Builder)

	t := template.Must(template.New(name).Parse(networksTemplates["external-network"]))
	err := t.Execute(out, map[string]string{
		"NetworkName": name,
	})
	if err != nil {
		return "", err
	}

	return out.String(), nil
}

func generateComposeFile(values TemplateValues) (string, error) {
	serverContainerName := "dank-server"
	if values.GoEnv == "beta" {
		serverContainerName += "-beta"
	}
	webContainerName := "dank-web-client"
	if values.GoEnv == "beta" {
		webContainerName += "-beta"
	}
	cdnContainerName := "dank-cdn"
	if values.GoEnv == "beta" {
		cdnContainerName += "-beta"
	}
	eventhubContainerName := "dank-eventhub"
	if values.GoEnv == "beta" {
		eventhubContainerName += "-beta"
	}
	ytdlContainerName := "dank-ytdl"
	if values.GoEnv == "beta" {
		ytdlContainerName += "-beta"
	}

	filesVolumeName := "dank-files"
	dbConfigVolumeName := "dank-db-config"
	dbDataVolumeName := "dank-db-data"

	var services []string
	var volumes []string
	var networks []string

	server, err := getService(ServiceValues{
		ServiceName:  serverContainerName,
		BuildContext: "./server",
		Dockerfile:   "Dockerfile",
		Ports: []ServicePortsValues{
			{
				OutPort:       values.ServerPort,
				ContainerPort: "3000",
			},
		},
		Environment: []ServiceEnvironmentValues{
			versionEnv,
			{
				Key:   "GO_ENV",
				Value: values.GoEnv,
			},
		},
		EnvFile: ".env.docker",
		Volumes: []ServiceVolumesValues{
			{
				VolumeName: filesVolumeName,
				MountPath:  "/app/.serve",
			},
		},
		Networks:  []string{values.NetworkName},
		DependsOn: []string{cdnContainerName, ytdlContainerName, eventhubContainerName},
	})
	if err != nil {
		return "", err
	}
	if slices.Contains(values.EnabledServices, ServiceServer) {
		services = append(services, server)
	}

	webClient, err := getService(ServiceValues{
		ServiceName:  webContainerName,
		BuildContext: "./web",
		Dockerfile:   "Dockerfile",
		Ports: []ServicePortsValues{
			{
				OutPort:       values.WebPort,
				ContainerPort: "3003",
			},
		},
		Environment: []ServiceEnvironmentValues{
			versionEnv,
			{
				Key:   "GO_ENV",
				Value: values.GoEnv,
			},
		},
		EnvFile:   ".env.docker",
		Volumes:   []ServiceVolumesValues{},
		Networks:  []string{values.NetworkName},
		DependsOn: []string{serverContainerName},
	})
	if err != nil {
		return "", err
	}
	if slices.Contains(values.EnabledServices, ServiceWeb) {
		services = append(services, webClient)
	}

	cdn, err := getService(ServiceValues{
		ServiceName:  cdnContainerName,
		BuildContext: "./server",
		Dockerfile:   "Dockerfile.cdn",
		Ports: []ServicePortsValues{
			{
				OutPort:       values.CdnPort,
				ContainerPort: "3001",
			},
		},
		Environment: []ServiceEnvironmentValues{
			versionEnv,
		},
		EnvFile: ".env.docker",
		Volumes: []ServiceVolumesValues{
			{
				VolumeName: filesVolumeName,
				MountPath:  "/app/.serve",
			},
		},
		Networks: []string{values.NetworkName},
	})
	if err != nil {
		return "", err
	}
	if slices.Contains(values.EnabledServices, ServiceCdn) {
		services = append(services, cdn)
	}

	eventhub, err := getService(ServiceValues{
		ServiceName:  eventhubContainerName,
		BuildContext: "./server",
		Dockerfile:   "Dockerfile.eventhub",
		Ports: []ServicePortsValues{
			{
				OutPort:       values.EventHubPort,
				ContainerPort: "3002",
			},
		},
		Environment: []ServiceEnvironmentValues{
			versionEnv,
		},
		EnvFile: ".env.docker",
		Volumes: []ServiceVolumesValues{
			{
				VolumeName: filesVolumeName,
				MountPath:  "/app/.serve",
			},
		},
		Networks: []string{values.NetworkName},
	})
	if err != nil {
		return "", err
	}
	if slices.Contains(values.EnabledServices, ServiceEventHub) {
		services = append(services, eventhub)
	}

	ytdl, err := getService(ServiceValues{
		ServiceName:    ytdlContainerName,
		BuildContext:   "./ytdl",
		Dockerfile:     "Dockerfile",
		ContainerImage: "",
		Ports: []ServicePortsValues{
			{
				OutPort:       values.YtDlPort,
				ContainerPort: "8000",
			},
		},
		Environment: []ServiceEnvironmentValues{},
		EnvFile:     ".env.docker",
		Volumes: []ServiceVolumesValues{
			{
				VolumeName: filesVolumeName,
				MountPath:  "/app/.serve",
			},
		},
		Networks: []string{
			values.NetworkName,
		},
		DependsOn: []string{},
	})
	if err != nil {
		return "", err
	}
	if slices.Contains(values.EnabledServices, ServiceYtDl) {
		services = append(services, ytdl)
	}

	db, err := getService(ServiceValues{
		ServiceName:    "dank-db",
		ContainerImage: values.MariaDbImage,
		Ports: []ServicePortsValues{
			{
				OutPort:       "3306",
				ContainerPort: "3306",
			},
		},
		Environment: []ServiceEnvironmentValues{
			{
				Key:   "MARIADB_ROOT_PASSWORD",
				Value: "previetcomrade",
			},
			{
				Key:   "MARIADB_DATABASE",
				Value: "dankabase",
			},
		},
		EnvFile: ".env.docker",
		Volumes: []ServiceVolumesValues{
			{
				VolumeName: dbConfigVolumeName,
				MountPath:  "/etc/mysql",
			},
			{
				VolumeName: dbDataVolumeName,
				MountPath:  "/var/lib/mysql",
			},
		},
		Networks: []string{values.NetworkName},
	})
	if err != nil {
		return "", err
	}
	if slices.Contains(values.EnabledServices, ServiceDb) {
		services = append(services, db)
	}

	if values.ExternalNetwork {
		network, err := getExternalNetwork(values.NetworkName)
		if err != nil {
			return "", err
		}
		networks = append(networks, network)
	} else {
		network, err := getInternalNetwork(values.NetworkName)
		if err != nil {
			return "", err
		}
		networks = append(networks, network)
	}

	uploads, err := getVolume(VolumeValues{
		VolumeName: filesVolumeName,
		LocalPath:  values.FilesPath,
	})
	if err != nil {
		return "", err
	}

	if slices.Contains(values.EnabledServices, ServiceCdn) || slices.Contains(values.EnabledServices, ServiceYtDl) {
		volumes = append(volumes, uploads)
	}

	if slices.Contains(values.EnabledServices, ServiceDb) {
		dbConfig, err := getVolume(VolumeValues{
			VolumeName: dbConfigVolumeName,
			LocalPath:  values.DbConfigPath,
		})
		if err != nil {
			return "", err
		}
		if dbConfig != "" {
			volumes = append(volumes, dbConfig)
		}
	}

	if slices.Contains(values.EnabledServices, ServiceDb) {
		dbData, err := getVolume(VolumeValues{
			VolumeName: dbDataVolumeName,
			LocalPath:  values.DbDataPath,
		})
		if err != nil {
			return "", err
		}
		if dbData != "" {
			volumes = append(volumes, dbData)
		}
	}

	cache, err := getService(ServiceValues{
		ServiceName:    "dank-cache",
		ContainerImage: values.RedisImage,
		Ports: []ServicePortsValues{
			{
				OutPort:       "6379",
				ContainerPort: "6379",
			},
		},
		Environment: []ServiceEnvironmentValues{},
		EnvFile:     ".env.docker",
		Volumes:     []ServiceVolumesValues{},
		Networks:    []string{values.NetworkName},
		Command:     "--requirepass previetcomrade",
	})
	if err != nil {
		return "", err
	}
	if slices.Contains(values.EnabledServices, ServiceCache) {
		services = append(services, cache)
	}

	out := new(strings.Builder)

	template := template.Must(template.New("the-thing").Parse(dockerComposeFullTemplate))
	err = template.Execute(out, map[string][]string{
		"Services": services,
		"Networks": networks,
		"Volumes":  volumes,
	})
	if err != nil {
		return "", err
	}

	return out.String(), nil
}

const (
	ServiceServer   string = "dank-server"
	ServiceWeb      string = "dank-web-client"
	ServiceCdn      string = "dank-cdn"
	ServiceEventHub string = "dank-eventhub"
	ServiceYtDl     string = "dank-ytdl"
	ServiceDb       string = "dank-db"
	ServiceCache    string = "dank-cache"
)

func generateAllComposeFile() error {
	content, err := generateComposeFile(TemplateValues{
		MariaDbImage:    mariadbImage,
		RedisImage:      redisImage,
		ServerPort:      "20250",
		WebPort:         "20253",
		CdnPort:         "20251",
		EventHubPort:    "20252",
		YtDlPort:        "20254",
		NetworkName:     "danknetwork",
		FilesPath:       "./.serve",
		DbConfigPath:    "./.db/etc",
		DbDataPath:      "./.db/var",
		ExternalNetwork: false,
		EnabledServices: []string{
			ServiceServer,
			ServiceWeb,
			ServiceCdn,
			ServiceEventHub,
			ServiceYtDl,
			ServiceDb,
			ServiceCache,
		},
	})
	if err != nil {
		return err
	}

	file, err := os.Create("../../docker-compose-all.yml")
	if err != nil {
		return err
	}

	_, err = file.WriteString(content)
	if err != nil {
		return err
	}

	return nil
}

func generateJustServicesComposeFile() error {
	content, err := generateComposeFile(TemplateValues{
		MariaDbImage:    mariadbImage,
		RedisImage:      redisImage,
		CdnPort:         "20251",
		EventHubPort:    "20252",
		YtDlPort:        "20254",
		NetworkName:     "danknetwork",
		FilesPath:       "./.serve",
		DbConfigPath:    "./.db/etc",
		DbDataPath:      "./.db/var",
		ExternalNetwork: false,
		EnabledServices: []string{
			ServiceCdn,
			ServiceEventHub,
			ServiceYtDl,
			ServiceDb,
			ServiceCache,
		},
	})
	if err != nil {
		return err
	}

	file, err := os.Create("../../docker-compose.yml")
	if err != nil {
		return err
	}

	_, err = file.WriteString(content)
	if err != nil {
		return err
	}

	return nil
}

func generateProdComposeFile(networkName, uploadsDir string) error {
	content, err := generateComposeFile(TemplateValues{
		ServerPort:      "20250",
		WebPort:         "20253",
		CdnPort:         "20251",
		EventHubPort:    "20252",
		YtDlPort:        "20254",
		NetworkName:     networkName,
		FilesPath:       uploadsDir,
		GoEnv:           "prod",
		ExternalNetwork: true,
		EnabledServices: []string{
			ServiceServer,
			ServiceWeb,
			ServiceCdn,
			ServiceEventHub,
			ServiceYtDl,
		},
	})
	if err != nil {
		return err
	}

	file, err := os.Create("../../docker-compose-prod.yml")
	if err != nil {
		return err
	}

	_, err = file.WriteString(content)
	if err != nil {
		return err
	}

	return nil
}

func generateBetaComposeFile(networkName, uploadsDir string) error {
	content, err := generateComposeFile(TemplateValues{
		ServerPort:      "20360",
		WebPort:         "20363",
		CdnPort:         "20361",
		EventHubPort:    "20362",
		YtDlPort:        "20364",
		NetworkName:     networkName,
		FilesPath:       uploadsDir,
		GoEnv:           "beta",
		ExternalNetwork: true,
		EnabledServices: []string{
			ServiceServer,
			ServiceWeb,
			ServiceCdn,
			ServiceEventHub,
			ServiceYtDl,
		},
	})
	if err != nil {
		return err
	}

	file, err := os.Create("../../docker-compose-beta.yml")
	if err != nil {
		return err
	}

	_, err = file.WriteString(content)
	if err != nil {
		return err
	}

	return nil
}

func main() {
	generationType := flag.String("type", "services", "Docker compose file type, pre-set types are `all`, `beta`, `prod`, `services`")
	networkName := flag.String("network", "danknetwork", "External network name, used with `type=prod|beta`")
	uploadsDir := flag.String("uploads", "./.serve", "Uploads path, used with `type=prod|beta`")

	flag.Parse()

	var err error
	switch *generationType {
	case "all":
		err = generateAllComposeFile()
	case "beta":
		err = generateBetaComposeFile(*networkName, *uploadsDir)
	case "prod":
		err = generateProdComposeFile(*networkName, *uploadsDir)
	default:
		err = generateJustServicesComposeFile()
	}
	if err != nil {
		panic(err)
	}
}
