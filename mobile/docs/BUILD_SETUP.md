# Build Setup Guide

## Android Keystore Setup

### Option 1: EAS Managed Keystore (Recommended)

EAS can manage your keystore automatically:

```bash
cd mobile
eas build:configure
eas credentials
```

Select "Android" → "Generate new keystore" to let EAS handle it.

### Option 2: Custom Keystore

1. Generate a keystore:
```bash
keytool -genkeypair -v -storetype PKCS12 -keystore localis-mobile.keystore -alias localis -keyalg RSA -keysize 2048 -validity 10000
```

2. Set up EAS credentials:
```bash
eas credentials
```

3. Select "Android" → "Upload existing keystore" and provide:
   - Keystore path: `./localis-mobile.keystore`
   - Keystore password
   - Key alias: `localis`
   - Key password

4. **IMPORTANT**: Never commit the keystore to git. Add to `.gitignore`:
```
*.keystore
*.jks
google-service-account.json
```

## iOS Certificates

### Option 1: EAS Managed (Recommended)

```bash
eas credentials
```

Select "iOS" → "Generate new certificates" to let EAS handle it.

### Option 2: Manual Setup

1. Create an Apple Developer account at https://developer.apple.com
2. Generate a Certificate Signing Request (CSR) via Keychain Access
3. Upload CSR to Apple Developer portal
4. Download the certificate
5. Upload to EAS:
```bash
eas credentials
```

## Building

### Development Build
```bash
eas build --profile development --platform android
```

### Preview Build (Internal Testing)
```bash
eas build --profile preview --platform android
```

### Production Build
```bash
eas build --profile production --platform android
```

### iOS Build
```bash
eas build --profile production --platform ios
```

## Environment Variables

Set API URL for builds:

```bash
# In eas.json, add env to build profiles:
"production": {
  "env": {
    "EXPO_PUBLIC_API_URL": "https://api.localis.com"
  }
}
```

## Testing Builds

### Android APK (Direct Install)
```bash
eas build --profile preview --platform android
```

### Android AAB (Play Store)
```bash
eas build --profile production --platform android
```

### iOS IPA (TestFlight)
```bash
eas build --profile production --platform ios
eas submit --platform ios
```

## Troubleshooting

### Build Fails

1. Check EAS logs:
```bash
eas build:list
eas build:view <build-id>
```

2. Verify app.json configuration:
```bash
npx expo config
```

3. Clear cache:
```bash
npx expo start --clear
```

### Keystore Issues

1. Verify keystore:
```bash
keytool -list -v -keystore localis-mobile.keystore
```

2. Reset credentials:
```bash
eas credentials --reset
```

## Production Checklist

- [ ] Android keystore generated/configured
- [ ] iOS certificates configured (if targeting iOS)
- [ ] API URL set to production
- [ ] App icon and splash screen assets
- [ ] Version number incremented
- [ ] Build tested on physical device
- [ ] App store metadata prepared