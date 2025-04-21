# imk (Invoice Maker)
Invoice maker is a tool I use to generate invoices for work. I manage a bunch of digital receipts, and wanted a way to streamline generating invoices.

I upload all my receipts to Google Drive, being sure to follow naming conventions, and then when I am ready to generate an invoice, I download the directory from Google Drive and run imk on it.

## Installation
```bash
go install github.com/phillip-england/imk@latest
```

## Usage
```bash
imk [TARGET_DIR] [OUT.txt FILE] [INVOICE NAME]
# imk ./pdfs out.txt "UTICA #2 EOM"
```

## Naming Conventions
imk expects your files to be named a certain way. Here is the format:
```bash
[DATE]-[VENDOR]-[COST]-[DESCRIPTION]-[CATEGORY]-[BUSINESS].pdf
# 040325-target-32.89-sharpies-office-utica.pdf
```

Thanks much 😀
