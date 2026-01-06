# Setting Up APT Repository Publishing

This guide explains how to set up the required GitHub secrets for publishing DEB packages to an APT repository on GitHub Pages.

## Prerequisites

1. A GitHub repository with GitHub Pages enabled
2. GPG key pair for signing the APT repository
3. GitHub repository secrets configured

## Step 1: Generate GPG Keys

Generate a GPG key pair for signing your APT repository:

```bash
# Generate a new GPG key
gpg --full-generate-key
```

When prompted:
- Choose: RSA and RSA
- Key size: 4096 bits
- Expiration: Choose based on your needs (recommended: 2 years)
- Enter your name and email

After generation, find your key ID:

```bash
gpg --list-secret-keys --keyid-format LONG
```

## Step 2: Export GPG Keys

Export your private and public keys:

```bash
# Replace YOUR_KEY_ID with your actual key ID
gpg --armor --export-secret-keys YOUR_KEY_ID > private.key
gpg --armor --export YOUR_KEY_ID > public.key
```

## Step 3: Add GitHub Secrets

Add the following secrets to your GitHub repository:

1. Go to your repository on GitHub
2. Navigate to Settings → Secrets and variables → Actions
3. Click "New repository secret"

Add these three secrets:

### APT_SIGNING_KEY
- Name: `APT_SIGNING_KEY`
- Value: Contents of `private.key` (entire file including BEGIN/END lines)

### APT_SIGNING_PUBKEY
- Name: `APT_SIGNING_PUBKEY`
- Value: Contents of `public.key` (entire file including BEGIN/END lines)

### APT_SIGNING_KEY_PASSPHRASE
- Name: `APT_SIGNING_KEY_PASSPHRASE`
- Value: The passphrase you used when creating the GPG key
- If you didn't set a passphrase, you can skip this secret and remove the `key_passphrase` line from the workflow

## Step 4: Enable GitHub Pages

1. Go to Settings → Pages
2. Set Source to "Deploy from a branch"
3. Select the `gh-pages` branch
4. Set folder to `/ (root)`
5. Click Save

## Step 5: Run the Workflow

The workflow will automatically run when you:
- Push a tag starting with `v` (e.g., `v1.0.0`)
- Manually trigger it from the Actions tab

### Manual Trigger:

1. Go to Actions tab
2. Click "Build and Publish DEB Packages"
3. Click "Run workflow"
4. Enter a version number (e.g., `1.0.0`)
5. Click "Run workflow"

### Tag-based Trigger:

```bash
git tag v1.0.0
git push origin v1.0.0
```

## Step 6: Verify Publication

After the workflow completes successfully:

1. Check the `gh-pages` branch for the `repo/` directory
2. Visit `https://YOUR_USERNAME.github.io/gomailserver/repo/` to see your APT repository
3. The public key should be at `https://YOUR_USERNAME.github.io/gomailserver/repo/public.key`

## Security Considerations

1. **Keep your private key secure**: Never commit it to the repository
2. **Rotate keys periodically**: Generate new keys every 1-2 years
3. **Use a strong passphrase**: Protect your private key with a strong passphrase
4. **Backup your keys**: Store them securely offline

## Troubleshooting

### Workflow fails with GPG error

- Verify the secrets are correctly pasted (no extra spaces or line breaks)
- Ensure the passphrase matches your key

### Packages not appearing in repository

- Check that the `gh-pages` branch exists
- Verify GitHub Pages is enabled and pointing to the correct branch
- Look at the workflow logs for any errors

### Users can't download packages

- Ensure GitHub Pages is publicly accessible
- Check that the public key is accessible at the expected URL
- Verify the repository structure in the `gh-pages` branch

## Cleanup

To securely delete the exported keys after adding them to GitHub:

```bash
shred -vfz -n 10 private.key public.key
```

## Additional Resources

- [apt-repo-action Documentation](https://github.com/smeinecke/apt-repo-action)
- [GPG Documentation](https://gnupg.org/documentation/)
- [GitHub Pages Documentation](https://docs.github.com/en/pages)
