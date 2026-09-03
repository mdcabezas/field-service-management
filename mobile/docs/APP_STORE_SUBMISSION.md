# App Store Submission Guide

## Google Play Store

### Prerequisites
- Google Play Developer Account ($25 one-time)
- App signing key configured

### Step 1: Create App Listing

1. Go to Google Play Console
2. Click "Create app"
3. Fill in:
   - App name: Localis Mobile
   - Default language: Spanish
   - App or game: App
   - Free or paid: Free

### Step 2: Store Listing

#### Short description (80 chars max)
"Aplicación móvil para gestión de visitas técnicas y reportes"

#### Full description (4000 chars max)
"Localis Mobile permite a los técnicos de campo:

• Gestionar visitas con checklist personalizados
• Capturar fotos con marca de tiempo y GPS
• Registrar checkpoints de ubicación
• Tomar mediciones técnicas
• Registrar uso de materiales
• Trabajar sin conexión y sincronizar después
• Generar reportes automáticos

Diseñada para dispositivos de gama baja con 2GB de RAM."

### Step 3: Graphics

#### App icon (512x512 PNG)
- Upload to Play Console

#### Feature graphic (1024x500 PNG)
- Required for store listing

#### Screenshots
- Phone: minimum 2, recommended 4-8
- Tablet: optional but recommended

### Step 4: Content Rating

Complete the IARC content rating questionnaire:
- Violence: None
- Language: None
- Controlled substances: None
- User interaction: None

### Step 5: Privacy Policy

Create and host a privacy policy URL:
- What data is collected
- How data is used
- Data retention policy
- Contact information

### Step 6: Data Safety

Declare data collection:
- Location: Collected
- Photos: Collected
- Device info: Collected
- Usage data: Not collected

### Step 7: Upload Build

1. Build AAB:
```bash
eas build --profile production --platform android
```

2. Download AAB from EAS

3. Upload to Play Console > App releases

4. Fill in release notes

### Step 8: Review and Publish

1. Review all sections
2. Submit for review
3. Wait for approval (typically 1-3 days)

---

## Apple App Store

### Prerequisites
- Apple Developer Account ($99/year)
- App Store Connect access

### Step 1: Create App

1. Go to App Store Connect
2. Click "My Apps" > "+"
3. Fill in:
   - Platform: iOS
   - Name: Localis Mobile
   - Primary Language: Spanish
   - Bundle ID: com.localis.mobile
   - SKU: localis-mobile

### Step 2: App Information

#### Name (30 chars max)
"Localis Mobile"

#### Subtitle (30 chars max)
"Gestión de visitas técnicas"

#### Category
Primary: Business
Secondary: Utilities

#### Content Rights
"No, this app does not contain third-party content"

### Step 3: Pricing and Availability

- Price: Free
- Availability: All countries

### Step 4: App Privacy

#### Data Types Collected
- Location
- Photos
- User Content
- Identifiers

#### Privacy Policy URL
Create and host a privacy policy

### Step 5: Version Information

#### What's New (170 chars max)
"Versión inicial de Localis Mobile con gestión de visitas, checklist, fotos y sincronización offline."

#### Screenshots
- iPhone 6.7": 1290 x 2796 or 1284 x 2778
- iPhone 6.5": 1242 x 2688 or 1284 x 2778
- iPhone 5.5": 1242 x 2208

#### App Preview (optional)
- 15-30 seconds
- Shows key functionality

### Step 6: Build

1. Build IPA:
```bash
eas build --profile production --platform ios
```

2. Submit to App Store Connect:
```bash
eas submit --platform ios
```

3. Select build in App Store Connect

### Step 7: Review

1. Complete all required fields
2. Submit for review
3. Wait for approval (typically 24-48 hours)

---

## Post-Submission

### Monitoring

1. Monitor crash reports
2. Check user reviews
3. Respond to feedback

### Updates

1. Increment version number
2. Update changelog
3. Submit new build

### Analytics

Consider adding:
- Firebase Analytics
- Crashlytics
- Performance monitoring

---

## Resources

- Google Play Console: https://play.google.com/console
- App Store Connect: https://appstoreconnect.apple.com
- Expo Submit: https://docs.expo.dev/submit/introduction/