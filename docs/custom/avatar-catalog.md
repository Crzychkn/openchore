# Curated profile pictures

Use **Me → Change picture** to preview and save a library image, generated
DiceBear picture or initials. Choosing a library image keeps the previous
legacy URL as a fallback if the asset is later unavailable. Saving a DiceBear
picture selects the generated image again. Saving initials clears both choices.

OpenChore ships three original geometric illustrations under the repository's
[MIT license](../../LICENSE). They are embedded in the Go server and work
without access to an external image service. The generated options still
use DiceBear.

## Administrator-managed images

A server operator can curate additional assets by setting
`OPENCHORE_AVATAR_DIR` for the API process to a directory backed by persistent
storage, such as `/data/personalization/avatars` in a container. Create the
directory before starting the API. Add up to 32 approved PNGs with lowercase
hyphenated names, for example `forest-fox.png`. The filename becomes the
display label and determines a stable ID such as `managed-forest-fox`.
Keep the filename when updating the artwork for an existing ID.

Prepare images with redistribution rights, at most 256 KiB and 512 × 512
pixels. Only regular files with valid names and PNG data enter the catalog.
The server loads and sanitizes images at startup. Add and review images in a
development instance before applying changes to a live catalog. Catalog images
are available to anyone who can reach the household's public profile views.

SQLite records the selected ID. Back up persistent managed image files along
with the database. The YAML seed export still exports legacy URLs and does
not carry curated selections or managed files. Back up the database to
preserve those selections.
