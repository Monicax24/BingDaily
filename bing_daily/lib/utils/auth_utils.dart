import 'package:firebase_auth/firebase_auth.dart';
import 'package:google_sign_in/google_sign_in.dart';

/// Handles Google Sign-In flow with Firebase Auth (v7+ compatible).
/// Returns the signed-in [User] or null if sign-in fails or domain invalid.
Future<User?> signInWithGoogle() async {
  try {
    final GoogleSignIn googleSignIn = GoogleSignIn.instance;

    // Initialize (call once, e.g., in initState or a service)
    await googleSignIn.initialize(
      clientId: '1098874852025-rkudo82gnti4un9764gq1h59uugasiop.apps.googleusercontent.com', // From Firebase (iOS/Web Client ID)
      serverClientId:
          '1098874852025-9cm4o3taotgntmu3gcokdkm8mknatulk.apps.googleusercontent.com', // From Firebase (Web Client ID for Android/server)
      hostedDomain:
          'binghamton.edu', // Attempt restriction (buggy; validate below)
    );

    // Optional: Attempt silent auth first
    await googleSignIn.attemptLightweightAuthentication();

    // Trigger interactive sign-in
    final GoogleSignInAccount googleUser = await googleSignIn.authenticate();


    // Get auth details
    final GoogleSignInAuthentication googleAuth = googleUser.authentication;

    // Create Firebase credential (idToken sufficient for basic auth)
    final OAuthCredential credential = GoogleAuthProvider.credential(
      idToken: googleAuth.idToken,
    );

    // Sign in to Firebase
    final UserCredential userCredential = await FirebaseAuth.instance
        .signInWithCredential(credential);

    // Workaround for hostedDomain bug: Validate email domain
    final String? email = userCredential.user?.email;
    if (email == null || !email.endsWith('@binghamton.edu')) {
      // Invalid domain: Sign out and return null
      await signOut();
      throw Exception('Invalid school domain');
    }

    return userCredential.user;
  } catch (e) {
    print('Google Sign-In Error: $e'); // Use logger in production
    return null;
  }
}

/// Signs out from Google and Firebase.
Future<void> signOut() async {
  await GoogleSignIn.instance.signOut();
  await FirebaseAuth.instance.signOut();
}
