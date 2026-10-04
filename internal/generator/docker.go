package generator

func writeDockerfile(dir string, data templateData) error {
	return render(dir, "Dockerfile", "templates/Dockerfile.tmpl", data)
}
