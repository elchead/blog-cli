# Blog-CLI ✍️

The utility is built for my blog workflow using Github pages with Hugo. It avoids a few repetitive steps I had to do before.
I like to use Obsidian for writing and the CLI allows me to publish my blog posts from Obisidian.

It should be easy to adjust the workflow, so feel free to reuse :)

## Features

`blog post <title>`:

- create new blog post skeleton (including metadata) in my writing app and (symbolically) links it to the corresponding `post` directory in my Github repo.
- after file creation, it opens the file in Obsidian

  `--book | -B`: create book note from template file

---

`blog draft <title>`:

- same as `post` without linkage in repo

  `--book | -B`: create book note from template file

---

`blog publish <title>`:

- use existing Obsidian article to create linkage in repo. Then locally render blog (`hugo serve`) and open preview in Browser

  `--book | -B`: create book note from template file

---

`blog preview`:

- render the blog (`hugo serve`) and open in Browser

---

### Still missing

- media (images) for posts still need to be manually added to the repo
- modify Metadata inside book template
- separate config file for setting paths

## Build

1. Modify the config paths inside `cmd/cmd.go`.

2. Change into the repo directory, then:

   `go build -o ./bin/blog ./cmd/cmd.go`

3. (Optional): create link in PATH (Mac):

   `` ln -s `pwd`/bin/blog /usr/local/bin/blog ``
