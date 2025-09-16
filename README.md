# issuehtml

# Поиск открытых issues в Go репозитории

go run main.go repo:golang/go is:open

# Поиск по ключевым словам

go run main.go "json marshal" language:go

# Поиск закрытых issues

go run main.go repo:google/gvisor is:closed
