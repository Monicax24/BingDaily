import 'package:firebase_auth/firebase_auth.dart';
import 'package:google_sign_in/google_sign_in.dart';

/// Handles Google Sign-In flow with Firebase Auth (v7+ compatible).
/// Returns the signed-in [User] or null if sign-in fails or domain is invalid.
Future<User?> signInWithGoogle() async {
  try {
    final GoogleSignIn googleSignIn = GoogleSignIn.instance;

    // Initialize Google Sign-In with client IDs and domain restriction
    await googleSignIn.initialize(
      clientId:
          '1098874852025-rkudo82gnti4un9764gq1h59uugasiop.apps.googleusercontent.com',
      serverClientId:
          '1098874852025-pm39nsve54gqsmbi3fn46m83kcnauk16.apps.googleusercontent.com',
      hostedDomain: 'binghamton.edu',
    );

    // Attempt silent authentication first
    await googleSignIn.attemptLightweightAuthentication();

    // Trigger interactive sign-in
    final GoogleSignInAccount googleUser = await googleSignIn.authenticate();

    // Get authentication details
    final GoogleSignInAuthentication googleAuth = googleUser.authentication;

    // Create Firebase credential using ID token
    final OAuthCredential credential = GoogleAuthProvider.credential(
      idToken: googleAuth.idToken,
    );

    // Sign in to Firebase with the credential
    final UserCredential userCredential = await FirebaseAuth.instance
        .signInWithCredential(credential);

    // Validate email domain for Binghamton restriction
    final String? email = userCredential.user?.email;
    if (email == null || !email.endsWith('@binghamton.edu')) {
      await signOut();
      throw Exception('Invalid school domain');
    }

    return userCredential.user;
  } catch (e) {
    print('Google Sign-In Error: $e'); // Use logger in production
    return null;
  }
}

/// Signs out from both Google and Firebase.
Future<void> signOut() async {
  await GoogleSignIn.instance.signOut();
  await FirebaseAuth.instance.signOut();
}
