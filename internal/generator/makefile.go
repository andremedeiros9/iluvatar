package generator

func writeMakefile(dir string, data templateData) error {
	return render(dir, "Makefile", "templates/Makefile.tmpl", data)
}
