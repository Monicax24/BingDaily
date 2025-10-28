import 'package:bing_daily/api/firebase_auth.dart';
import 'package:bing_daily/constants.dart';
import 'package:bing_daily/models/app_user.dart';
import 'package:bing_daily/providers/auth_provider.dart';
import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';

// Login screen with Google Sign-In button for user authentication.
class LoginScreen extends ConsumerWidget {
  const LoginScreen({super.key});

  // Triggers Google Sign-In, fetches user data, and navigates to username setup.
  Future<void> _handleSignIn(BuildContext context, WidgetRef ref) async {
    final user = await signInWithGoogle();
    if (user != null) {
      final AppUser? appUser = await fetchUser();
      if (appUser != null) {
        ref.read(userNotifierProvider.notifier).setUser(appUser);
        print("set user successfully");
        context.go('/username');
      } else {
        ScaffoldMessenger.of(context).showSnackBar(
          const SnackBar(content: Text('Failed to fetch user data')),
        );
      }
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
            const Padding(padding: EdgeInsets.only(top: 200)),
            const Text(
              "Welcome to BingDaily!",
              style: TextStyle(fontSize: 35.0),
            ),
            const Padding(padding: EdgeInsets.only(top: 100)),
            ElevatedButton.icon(
              onPressed: () => _handleSignIn(context, ref),
              style: ElevatedButton.styleFrom(
                shape: RoundedRectangleBorder(
                  borderRadius: BorderRadius.circular(0.0),
                ),
                backgroundColor: bingWhite,
                foregroundColor: Colors.black,
              ),
              icon: Image.asset('assets/images/google_icon.jpg', height: 30),
              iconAlignment: IconAlignment.start,
              label: const Text(
                loginButtonText,
                style: TextStyle(fontFamily: loginButtonFont),
              ),
            ),
          ],
        ),
      ),
      backgroundColor: bingWhite,
    );
  }
}
