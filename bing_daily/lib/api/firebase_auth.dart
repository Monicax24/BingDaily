// Handles Firebase Authentication operations for the Bing Daily app.
import 'dart:convert';

import 'package:bing_daily/models/app_user.dart';
import 'package:firebase_auth/firebase_auth.dart';
import 'package:google_sign_in/google_sign_in.dart';
import 'package:http/http.dart' as http;

// Base URL for backend API
const String _baseUrl = 'http://18.234.86.119:8080';

// Handles Google Sign-In flow with Firebase Auth (v7+ compatible).
// Returns the signed-in [User] or null if sign-in fails or domain is invalid.
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

// Fetches user data from backend and returns AppUser.
Future<AppUser?> fetchUser() async {
  try {
    final headers = await _getHeaders();
    final userId = FirebaseAuth.instance.currentUser!.uid;
    final response = await http.get(
      Uri.parse('$_baseUrl/user/$userId'),
      headers: headers,
    );
    if (response.statusCode != 200) {
      throw Exception('Failed to fetch user: ${response.body}');
    }
    return AppUser.fromMap(jsonDecode(response.body)['data']);
  } catch (e) {
    print('Fetch User Error: $e'); // Use logger in production
    return null;
  }
}

// Retrieves Firebase ID token for authenticated requests.
Future<Map<String, String>> _getHeaders() async {
  final token = await _getIdToken();
  if (token == null) throw Exception('No authentication token available');
  return {'Authorization': 'Bearer $token'};
}

// Gets Firebase ID token for the current user.
Future<String?> _getIdToken() async {
  final user = FirebaseAuth.instance.currentUser;
  if (user == null) return null;
  try {
    return await user.getIdToken(true);
  } catch (e) {
    return null;
  }
}

// Signs out the user from Google and Firebase.
Future<void> signOut() async {
  await GoogleSignIn.instance.signOut();
  await FirebaseAuth.instance.signOut();
}
