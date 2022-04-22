# euruswallet

Crypto wallet example.
### Features
- Crypto wallet example.

You need to declare variable in `android/app/src/main/AndroidManifest.xml`
after adding this line, you can enable http connection on webview

```xml
<manifest ...>
    <uses-permission android:name="android.permission.INTERNET" />
    <application
        ...
        android:usesCleartextTraffic="true"
        ...>
        ...
    </application>
</manifest>
```

You need to declare variable in `android/app/src/main/AndroidManifest.xml`
after adding this line, you can use chatbot page

```xml
 <manifest ...>
   <!-- ... other tags -->
   <application ...>
  <activity android:name="com.livechatinc.inappchat.ChatWindowActivity" android:configChanges="orientation|screenSize" />
   </application>
 </manifest>

```

You need to declare add intent filters in `android/app/src/main/AndroidManifest.xml`:
You can use this format to launch app from deep link eurus://pay on android app

```xml
<manifest ...>
  <!-- ... other tags -->
  <application ...>
    <activity ...>
      <!-- ... other tags -->

      <!-- Deep Links -->
      <intent-filter>
        <action android:name="android.intent.action.VIEW" />
        <category android:name="android.intent.category.DEFAULT" />
        <category android:name="android.intent.category.BROWSABLE" />
        <!-- Accepts URIs that begin with YOUR_SCHEME://YOUR_HOST -->
        <data
          android:scheme="eurus"
          android:host="pay" />
      </intent-filter>

    </activity>
  </application>
</manifest>
```

You can use this format to launch app from deep link eurus://pay on ios app
For **Custom URL schemes** you need to declare the scheme in
`ios/Runner/Info.plist` (or through Xcode's Target Info editor,
under URL Types):

```xml
<?xml ...>
<!-- ... other tags -->
<plist>
<dict>
  <!-- ... other tags -->
  <key>CFBundleURLTypes</key>
  <array>
    <dict>
      <key>CFBundleTypeRole</key>
      <string>Editor</string>
      <key>CFBundleURLName</key>
      <string>pay</string>
      <key>CFBundleURLSchemes</key>
      <array>
        <string>eurus</string>
      </array>
    </dict>
  </array>
  <!-- ... other tags -->
</dict>
</plist>
```