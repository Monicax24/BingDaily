import 'package:bing_daily/api/auth_api.dart';
import 'package:bing_daily/constants.dart';
import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';

/// Login screen with Google Sign-In button for user authentication.
class LoginScreen extends StatelessWidget {
  const LoginScreen({super.key});

  /// Triggers Google Sign-In and navigates to username setup on success.
  void _handleSignIn(BuildContext context) async {
    final result = await signInWithGoogle();
    if (result != null) {
      context.go('/username');
    } else {
      ScaffoldMessenger.of(
        context,
      ).showSnackBar(const SnackBar(content: Text('Sign-in failed')));
    }
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(title: const Text(appTitle), backgroundColor: bingGreen),
      body: Center(
        child: Column(
          mainAxisAlignment: MainAxisAlignment.center,
          children: [
            ElevatedButton(
              onPressed: () => _handleSignIn(context),
              style: ElevatedButton.styleFrom(
                backgroundColor: bingAccent,
                foregroundColor: bingWhite,
              ),
              child: const Text(loginButtonText),
            ),
          ],
        ),
      ),
      backgroundColor: bingWhite,
    );
  }
}
