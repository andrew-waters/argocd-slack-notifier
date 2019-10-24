# Example helm template

Add your variables to `values.yaml` (these will be set as environment variables on the `Job`).

You can then run the following to generate your manifest file or you can deploy using helm:

```bash
helm template . > manifests/install.yaml
```

To deploy the generated manifest (outside of helm) make sure you run `kubectl create ./manifests/install.yaml` instead of `apply`.
