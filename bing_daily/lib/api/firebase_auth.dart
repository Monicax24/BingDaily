import 'package:bing_daily/models/app_user.dart';
import 'package:bing_daily/providers/auth_provider.dart';
import 'package:firebase_auth/firebase_auth.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:google_sign_in/google_sign_in.dart';
import 'package:http/http.dart' as http;

/// Handles Google Sign-In flow with Firebase Auth (v7+ compatible).
/// Returns the signed-in [User] or null if sign-in fails or domain is invalid.
final String _baseUrl = 'http://18.234.86.119:8080';
Future<User?> signInWithGoogle(WidgetRef ref) async {
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

    fetchUser(ref);

    return userCredential.user;
  } catch (e) {
    print('Google Sign-In Error: $e'); // Use logger in production
    return null;
  }
}



Future<String?> _getIdToken() async {
  final user = FirebaseAuth.instance.currentUser;
  if (user == null) return null;
  try {
    return await user.getIdToken(true);
  } catch (e) {
    return null;
  }
}

Future<Map<String, String>> _getHeaders() async {
  final token = await _getIdToken();
  if (token == null) throw Exception('No authentication token available');
  return {'Authorization': 'Bearer $token'};
}

Future<void> fetchUser(WidgetRef ref) async {
  final headers = await _getHeaders();
  final userId = FirebaseAuth.instance.currentUser!.uid;
  final response = await http.get(
    Uri.parse('$_baseUrl/user/$userId'),
    headers: headers,
  );
  if (response.statusCode != 200) {
    throw Exception('Failed to fetch user: ${response.body}');
  }
  ref
      .watch(userNotifierProvider.notifier)
      .setUser(AppUser.fromJson(response.body));

  print('User data: ${response.body}');
}

Future<void> signOut() async {
  await GoogleSignIn.instance.signOut();
  await FirebaseAuth.instance.signOut();
}
