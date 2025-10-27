import 'package:bing_daily/api/firebase_auth.dart';
import 'package:bing_daily/constants.dart';
import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';

/// Login screen with Google Sign-In button for user authentication.
class LoginScreen extends ConsumerWidget {
  const LoginScreen({super.key});

  /// Triggers Google Sign-In and navigates to username setup on success.
  void _handleSignIn(BuildContext context, WidgetRef ref) async {
    final result = await signInWithGoogle(ref);
    if (result != null) {
      context.go('/username');
    } else {
      ScaffoldMessenger.of(
        context,
      ).showSnackBar(const SnackBar(content: Text('Sign-in failed')));
    }
  }

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    return Scaffold(
      appBar: AppBar(title: const Text(appTitle), backgroundColor: bingGreen),
      body: Center(
        child: Column(
          mainAxisAlignment: MainAxisAlignment.start,
          children: [
            Padding(
              padding: EdgeInsets.only(
                top: 200,
              ),
            ),
            Text(
              "Welcome to BingDaily!",
              style: TextStyle(
                fontSize: 35.0,
              ),
            ),
            Padding(
              padding: EdgeInsets.only(
                top: 100,
              ),
            ),
            ElevatedButton.icon(
              onPressed: () => _handleSignIn(context, ref),
              style: ElevatedButton.styleFrom(
                shape: RoundedRectangleBorder(
                  borderRadius: BorderRadius.circular(0.0), // For a perfect rectangle
                ),
                backgroundColor: bingWhite,
                foregroundColor: Colors.black,
              ),
              icon: Image.asset(
                'assets/images/google_icon.jpg',
                height: 30,
              ),
              iconAlignment: IconAlignment.start,
              label: const Text(
                loginButtonText,
                style: TextStyle(
                  fontFamily: loginButtonFont,
                )
              ),
            ),
          ],
        ),
      ),
      backgroundColor: bingWhite,
    );
  }
}
