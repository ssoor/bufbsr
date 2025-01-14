# `codegenerator`

`codegenerator` is an implementation of remote code generation for custom
plugins to integrate with the `buf` CLI.

It implements the `buf.alpha.registry.v1alpha1.CodeGenerationService`
service interface.

It expects a locally executable registry to be available at
`CODEGENERATOR_REGISTRY_PATH`.

The format of this path must be `<owner>/<plugin>/<version>/<plugin>`

* `plugin` is the name of the binary (e.g. `protoc-gen-doc`).
* `version` must be of the from `v\d+\.\d+\.\d+` 

Assuming you host this service at `codegenerator.build` you can reference your
plugins in `buf.gen.yaml` as follows:

```yaml
version: v2

plugins:
 - remote: codegenerator.build/<owner>/<plugin>:<version>
   out: generated
```
