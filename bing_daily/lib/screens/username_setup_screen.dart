import 'package:bing_daily/constants.dart';
import 'package:flutter/material.dart';

/// Placeholder screen for setting up a username after Google Sign-In.
class UsernameSetupScreen extends StatelessWidget {
  const UsernameSetupScreen({super.key});

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: const Text(usernameSetupTitle),
        backgroundColor: bingGreen,
      ),
      body: const Center(
        child: Text(
          'Username setup coming soon!',
          style: TextStyle(color: bingGreen, fontSize: 20),
        ),
      ),
      backgroundColor: bingWhite,
    );
  }
}
