## v0.6.0 (Sat, 03 Oct 2026 11:49:31 UTC)
- fix: follow every pointer when building type metadata.
- fix: cache function names per function, not per type.
- fix: return nil for a negative FieldByIndex.
- fix: return nil when StructValue has no struct to read.
- fix: leave unsettable fields unchanged in NewIfNil.
- fix: publish one metadata when the type cache misses.
- fix: return copies of Index and Fields.
- fix: omit the dot when PackageAndName has an empty part.
- docs: correct when ErrInvField is returned.
- refactor: drop the unused Field.exported flag.
- test: match misnamed cases to the value they build.
- refactor: keep ErrTagSyntax with the other sentinels.
- test: order Field.IsExported with the method.
- test: order FieldValue.StructValue with the method.
- test: name StructValue tests after the method.
- test: name the NewField result have.
- test: name the NewFieldValue result have.
- test: name the Fields result have.
- test: name the nil StructValue result have.
- test: name error cases after the failing input.
- test: separate Then subjects in NewIfNil and StructValue.
- test: prepare the cached ReflectValue in Given.
- chore: bump testing to v0.56.0 and testkit to v0.15.0.
- docs: document the mirror import path and Go 1.26.
- refactor: use reflect.Pointer and reference TStruct.fStr.
- fix: allocate one pointer level in NewIfNil.
- fix: copy the field index out of reflect.
- fix: reject an invalid value before reading its type.
- test: prepare helper arguments in Given.
- test: assert the type indirect returns.
- test: pass the expected type before the result.
- test: drop slashes from benchmark names.
- test: shorten Tag.IsZero subtest names.
- test: distinguish the anonymous funcPkg cases.
- test: shorten the index mutation subtest name.
- test: inline the one-line ptr helper.
- test: assert NewTypeMetadata directly.
- docs: add an example for ReflectType.
- docs: add an example for ParseTags.
- docs: add an example for FieldValue.NewIfNil.
- docs: add an example for StructValue.NewIfNil.
- docs: format the ParseTags nolint directive.
- docs: describe which tag parts IsZero checks.
- docs: correct Field, IndirectType, and TypeMetadata.
- docs: drop Field godoc that restates the signature.
- docs: describe FieldByName benchmarks as map lookups.
- docs: describe what NewIfNil allocates.
- docs: finish the no-dot comment.
- docs: drop Tag godoc that restates the signature.
- docs: drop Metadata godoc that restates the signature.
- docs: drop FieldValue godoc that restates the signature.
- docs: drop StructValue godoc that restates the signature.
- fix!: return tag parse errors from Field.Tag.
- docs: correct Reflect, NewIfNil, and Get in the README.

## v0.5.0 (Mon, 08 Jun 2026 07:20:31 UTC)
- chore: bump deps, fix godoc errors, and clean up tests.
- perf: add benchmarks.
- perf: add field name index for O(1) FieldByName lookups.
- docs: sync README examples with gmdoceg markers.
- chore: update copyright year to 2026 and bump testkit to v0.3.1.

## v0.4.0 (Fri, 01 May 2026 19:27:03 UTC)
- chore: Update to Go 1.26 and update dependencies.

## v0.3.2 (Sun, 26 Apr 2026 13:34:10 UTC)
- chore: Update dependencies.

## v0.3.1 (Thu, 12 Feb 2026 13:51:08 UTC)
- style: Remove unneeded test case.
- chore: Update dependencies.

## v0.3.0 (Wed, 05 Nov 2025 08:39:08 UTC)
- chore: Update dependencies.
- Extract type name and import path when possible.

## v0.2.0 (Sun, 28 Sep 2025 19:09:08 UTC)
- style: Add missing SPDX lines to examples_test.go file.
- feat!: Rename `mirror.MetadataFor` and `mirror.TypeMetadata` to `mirror.Reflect` and `mirror.ReflectType`.

## v0.1.0 (Sat, 27 Sep 2025 15:40:18 UTC)
- feat: Initial commit.
- doc: Add documentation and examples.
- test: Add GitHub workflows.

