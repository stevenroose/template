# template
A simple command line tool for filling file templates based on Go's text/template

Make templates from any text file, f.e. a Markdown file `template.md`:

```markdown
# Welcome
My name is {{.name}}, contact me at {{.email}}!
```

Quickly and easily fill in the template variables using a single command:

`$ template -i template.md -o final.md -v name="Alice" -v email="alice@example.org`

Template variables can also be provided as a YAML or JSON file:

```
$ template -i template.yaml -o final.yaml -y data.yaml
$ template -i template.yaml -o final.yaml -j data.json
```