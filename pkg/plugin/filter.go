package plugin

type Filter func(Plugin) bool

func NameFilter(names ...string) Filter {
	mp := map[string]struct{}{}
	for _, name := range names {
		mp[name] = struct{}{}
	}
	return func(plugin Plugin) bool {
		_, ok := mp[plugin.Name()]
		return !ok
	}
}
