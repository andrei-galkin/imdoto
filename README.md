# imdoto

`imdoto` is a small Go CLI that searches image engines (Google, Bing, Yandex) and downloads matching images to a local folder.

Important: respect copyright and the terms of the image providers when using this tool.

**Status**: buildable and tested locally. Tests are included for the search packages.

## Quick start

- Build the CLI (from the repository root):

```powershell
cd C:\dev\imdoto
go build ./imdoto
```

- Run the produced executable (`imdoto.exe`) with flags:

```powershell
.\imdoto.exe -engine "google" -term "apple seed" -folder "images" -limit 100 -type "jpeg"
```

Flags:

- `-engine` : `google`, `bing`, or `yandex` (default: `yandex`)
- `-term`   : search term (words with spaces should be quoted)
- `-folder` : destination folder name (default: `img`)
- `-limit`  : maximum number of images to download (default: `75`)
- `-type`   : image file type filter (e.g. `jpeg`, `png`, or `*` for any)

## Run tests

```powershell
cd C:\dev\imdoto
go test ./...
```

## Formatting & vet

```powershell
go fmt ./...
go vet ./...
```

## Notes

- Module file: `go.mod` has module path `github.com/andrei-galkin/imdoto`.
- A `.gitignore` exists to ignore build artifacts and editor files.

## Contributing

- Open issues or PRs. Consider adding CI (GitHub Actions) to run `go test` and a linter such as `golangci-lint`.

## License

See `Licence.txt`.
