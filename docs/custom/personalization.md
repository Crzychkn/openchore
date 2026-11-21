# Future personalization

Reuse the Go API, SQLite user records, React settings and design primitives
described in the [architecture audit](architecture.md). Model both features
as selections from curated libraries. Household members can choose library
entries for themselves. The existing avatar options remain available,
including DiceBear, initials and legacy avatar URLs already stored or
supported by existing clients and imports. Arbitrary user file uploads
are outside this design.

## Asset catalog and persistent storage

- Give each entry a stable ID, kind (avatar or background), display label,
  attribution and availability. Supply bundled defaults under stable IDs
  with redistribution rights. Record the creator, source and license.
  Start with a small set of recognizable avatar images and backgrounds
  with enough negative space and contrast for the existing skins.
- Use Go to list catalog metadata, resolve IDs and serve local assets.
  Administrator-managed images belong under a dedicated directory on
  the persistent data volume, such as `/data/personalization`. Keep
  bundled resources in the application image. Use server-generated storage
  names and validate image type, size and decoded dimensions before
  committing an administrative import. Serve approved images via catalog
  routes separate from chore photos. Keep filesystem paths out of the API.
- Use SQLite for user selection IDs and, as needed, catalog metadata.
  Store image files on disk. Use kind-specific ID validation and keep
  asset revisions and retirements predictable for existing selections.
  Write and validate replacement files before publishing a reference,
  and preserve files until they are no longer referenced. Back up and
  restore the SQLite database and managed asset directory together.

## Profile picture selection

Use the existing Me sheet and shared Avatar component. Add a keyboard
accessible grid with preview and explicit save alongside the existing
picture chooser. A nullable catalog ID on the user gives existing rows
their current behavior. Prefer the selected catalog image, then the
existing `avatar_url`, then the initial fallback. When a user saves a
legacy avatar option, clear the catalog ID atomically so the change is
visible. A return to initials clears both. Preserve avatar behavior for
older clients and avoid overwriting a saved picture on unrelated admin edits.

Catalog avatars can appear in the household-visible profile picker, personal
header, Me sheet, People list and ambient avatar fallback. Preserve the
ambient pet presentation where it already has priority. Keep the public
projection limited to display data.

## Background selection

Add a curated background picker to personal appearance settings. Store a
nullable background ID on the user. Preview the selected image with the
current skin and provide a reset to the skin default. Keep the selection in
authenticated self-profile data. Apply the image behind KidApp content at its
personal `SkinScope`, retaining readable content surfaces and focus styles.
Use the skin default on missing assets or image errors. Shared House screens
keep their existing appearance, and switching users must drop the previous
background.

Use additive migrations, targeted updates for each preference, and the
existing authenticated self-service authorization pattern. Administrative
catalog changes require administrative authorization. Decide how selected
IDs are represented in YAML import/export and how imports behave when an
administrator-managed asset is missing.

## M2 proposal

Implement the curated avatar picker as the first vertical slice. Finalize
asset rights and catalog IDs, agree on the administrator-managed storage
and import lifecycle, then add the nullable avatar selection, API validation
and React picker. Verify existing DiceBear and imported URLs, selection and
reset, pre-M2 database upgrades, permissions, rendering across profile
surfaces, persistence across restart and operation without browser access
to third-party resources for curated images. Treat backgrounds as the next
slice, using the same catalog and storage design.
