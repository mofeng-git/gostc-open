package todo

import "server/repository"

// 修复一些线上数据错误的bug
func fix() {
	db, _, _ := repository.Get("")
	hosts, _ := db.GostClientHost.Where(db.GostClientHost.Or(
		db.GostClientHost.CustomDomain.IsNotNull(),
	).Or(
		db.GostClientHost.CustomDomain.Neq(""),
	)).Find()
	var effectiveDomainMap = make(map[string]bool)
	for _, host := range hosts {
		effectiveDomainMap[host.CustomDomain] = true
	}

	domains, _ := db.GostClientHostDomain.Find()
	var invalidDomains []string
	for _, domain := range domains {
		if effectiveDomainMap[domain.Domain] {
			continue
		}
		invalidDomains = append(invalidDomains, domain.Domain)
	}
	_, _ = db.GostClientHostDomain.Where(
		db.GostClientHostDomain.Domain.In(invalidDomains...),
	).Delete()
}
