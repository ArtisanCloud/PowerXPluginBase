# Quickstart — Managing Capabilities & Schemas

1. **Create Capability Descriptor**  
   - Add `contracts/capabilities/<domain>.<resource>.<action>.yaml` with ID/type/version/RBAC mapping.
   - Reference existing schema IDs or create new ones under `contracts/schema/input|output`.

2. **Author JSON Schemas**  
   - Use draft-07; place input/output definitions in the corresponding directories.  
   - Ensure schema `$id` matches file path for tooling to resolve.

3. **Reference in Manifest**  
   - Under `manifest.yaml`, list capability IDs in `capabilities.provides/consumes` with version hints.  
   - Keep detailed metadata out of manifest—only ID/type/version.

4. **Map RBAC Permissions**  
   - Update capability YAML `rbac` block; run `make check-capability` to ensure resource/action pairs exist in manifest RBAC contracts.

5. **Run Compatibility Checks**  
   - Execute `make check-compat` to diff schemas/capabilities against baseline versions.  
   - Review generated reports under `build/compat/` (or configured path).

6. **Update Documentation**  
   - Document capability in `docs/lifecycle/capabilities.md` (to be introduced) and sync integration docs with `make sync-lifecycle-docs`.

7. **Package & Release**  
   - `make verify-manifest` + `make package-pxp` ensures capability/schema assets ship in `.pxp`.  
   - Attach compatibility report and schema diffs to release notes under `docs/releases/`.
