# Deployment Checklist

## Android Deployment

### Pre-Deployment
- [x] All unit tests passing (66/66)
- [x] All backend tests passing (30/30)
- [x] TypeScript compilation clean
- [x] Bundle size < 15MB (4.2MB actual)
- [x] Performance optimizations applied
- [x] EAS configuration complete

### Build
- [ ] Generate release keystore
- [ ] Configure EAS credentials
- [ ] Run `eas build --profile production --platform android`
- [ ] Test APK on physical device
- [ ] Verify all features work offline

### Release
- [ ] Upload AAB to Google Play Console
- [ ] Complete store listing
- [ ] Set up internal testing track
- [ ] Distribute to testers

## iOS Deployment (Future)

### Requirements
- Apple Developer Account ($99/year)
- Mac with Xcode installed
- iOS device for testing

### Build
- [ ] Configure iOS certificates
- [ ] Run `eas build --profile production --platform ios`
- [ ] Test IPA on physical device
- [ ] Upload to App Store Connect

## Post-Deployment
- [ ] Monitor crash reports
- [ ] Collect user feedback
- [ ] Plan next iteration

## Known Limitations
1. **Timestamp watermark**: Requires native module (react-native-text-pipe)
2. **iOS build**: Requires Apple Developer account
3. **App store submission**: Requires store accounts and approval process

## Support
- Backend API: `http://localhost:8081`
- Documentation: `mobile/docs/`
- Issues: Report via project channels