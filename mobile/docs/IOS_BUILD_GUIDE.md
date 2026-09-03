# iOS Build Guide

## Prerequisites

### Required
- Apple Developer Account ($99/year)
- Mac with macOS 12.0+
- Xcode 14.0+
- CocoaPods installed

### Optional
- Fastlane for automation
- TestFlight for beta testing

## Step 1: Apple Developer Setup

1. Create Apple Developer account at https://developer.apple.com
2. Create an App ID in Certificates, Identifiers & Profiles
3. Register your iOS device for testing
4. Create a Distribution Certificate

## Step 2: EAS Configuration

### Configure iOS in eas.json
```json
{
  "build": {
    "production": {
      "ios": {
        "autoIncrement": true,
        "simulator": false
      }
    }
  }
}
```

### Set Apple Team ID
```bash
eas credentials
```

Select:
1. iOS
2. Apple Team
3. Enter your Team ID

## Step 3: Configure app.json

Add iOS-specific configuration:
```json
{
  "expo": {
    "ios": {
      "bundleIdentifier": "com.localis.mobile",
      "supportsTablet": true,
      "infoPlist": {
        "NSCameraUsageDescription": "Permitir acceso a la cámara para fotos de visitas",
        "NSLocationWhenInUseUsageDescription": "Permitir acceso a la ubicación para checkpoints GPS",
        "NSPhotoLibraryUsageDescription": "Permitir acceso a la galería para fotos de visitas"
      }
    }
  }
}
```

## Step 4: Build for iOS

### Development Build
```bash
eas build --profile development --platform ios
```

### Preview Build (TestFlight)
```bash
eas build --profile preview --platform ios
```

### Production Build
```bash
eas build --profile production --platform ios
```

## Step 5: Test with TestFlight

1. Upload build to App Store Connect:
```bash
eas submit --platform ios
```

2. Install TestFlight app on iOS device
3. Accept invitation via email
4. Install and test the build

## Step 6: App Store Submission

### Prepare App Store Listing
1. Screenshots (6.5" and 5.5" iPhones)
2. App description
3. Keywords
4. Support URL
5. Privacy Policy URL

### Submit for Review
1. Log in to App Store Connect
2. Select your app
3. Complete all required fields
4. Submit for review

## Common Issues

### Certificate Issues
```bash
# Reset credentials
eas credentials --reset

# List certificates
security find-identity -v -p codesigning
```

### Build Failures
```bash
# Clear cache
cd ios && pod deinstall && pod install
cd ..

# Reset EAS build
eas build --platform ios --clear
```

### Provisioning Profile Issues
1. Check bundle identifier matches
2. Verify device is registered
3. Ensure certificate is valid

## Environment Variables

Set in eas.json:
```json
{
  "build": {
    "production": {
      "env": {
        "EXPO_PUBLIC_API_URL": "https://api.localis.com"
      }
    }
  }
}
```

## Testing Checklist

- [ ] Login flow works
- [ ] Visit list loads
- [ ] Checklist functionality
- [ ] Photo capture works
- [ ] GPS checkpoints work
- [ ] Offline sync works
- [ ] Memory usage < 150MB
- [ ] No crashes on startup

## Resources

- Expo iOS Setup: https://docs.expo.dev/deploy/build-project/ios-specific-app-configuration/
- Apple Developer: https://developer.apple.com/documentation/
- TestFlight: https://developer.apple.com/testflight/