package internal

func getExternals(repo *Repository) []string {
	externals := make([]string, len(repo.DevDependencies))
	i := 0
	for external := range repo.DevDependencies {
		externals[i] = external
		i++
	}
	return externals
}
