import 'package:bing_daily/api/firebase_auth.dart';
import 'package:bing_daily/constants.dart';
import 'package:bing_daily/providers/auth_provider.dart';
import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';

/// Account page displaying user profile info and sign-out option.
class AccountPage extends ConsumerWidget {
  const AccountPage({super.key});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final user = ref.watch(userNotifierProvider);
    final username = user?.username ?? 'Guest';
    final email = user?.email ?? 'No email';
    final joinDate = user?.joinDate ?? 'Unknown';

    // Placeholder avatar with first letter of username
    final String initial = username.isNotEmpty
        ? username[0].toUpperCase()
        : 'G';

    return Scaffold(
      appBar: AppBar(
        title: const Text(
          accountLabel,
          style: TextStyle(color: bingWhite, fontWeight: FontWeight.bold),
        ),
        backgroundColor: bingGreen,
      ),
      backgroundColor: bingWhite,
      body: Center(
        child: Column(
          mainAxisAlignment: MainAxisAlignment.start,
          children: [
            const SizedBox(height: 40),
            // Profile picture placeholder
            CircleAvatar(
              radius: 50,
              backgroundColor: bingGreen,
              child: Text(
                initial,
                style: const TextStyle(
                  fontSize: 48,
                  color: bingWhite,
                  fontWeight: FontWeight.bold,
                ),
              ),
            ),
            const SizedBox(height: 24),
            // Username
            Text(
              username,
              style: const TextStyle(
                fontSize: 24,
                fontWeight: FontWeight.bold,
                color: bingGreen,
              ),
            ),
            const SizedBox(height: 8),
            // Email
            Text(
              email,
              style: TextStyle(fontSize: 16, color: Colors.grey[700]),
            ),
            const SizedBox(height: 8),
            // Join date
            Text(
              'Member since ${joinDate.substring(0, 10)}',
              style: TextStyle(fontSize: 14, color: Colors.grey[600]),
            ),
            const SizedBox(height: 80),
            // Sign out button
            GestureDetector(
              onTap: () async => await signOut(),
              child: Container(
                padding: const EdgeInsets.symmetric(
                  vertical: 12,
                  horizontal: 30,
                ),
                decoration: BoxDecoration(
                  color: Colors.grey[200],
                  borderRadius: BorderRadius.circular(12),
                  border: Border.all(color: Colors.grey[300]!),
                ),
                child: const Text(
                  'Sign Out',
                  style: TextStyle(fontSize: 12, color: Colors.blue),
                ),
              ),
            ),
          ],
        ),
      ),
    );
  }
}
