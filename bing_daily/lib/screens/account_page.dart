import 'package:bing_daily/api/auth_api.dart';
import 'package:flutter/material.dart';

class AccountPage extends StatelessWidget {
  const AccountPage({super.key});

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(title: const Text('Bing Daily')),
      body: Center(
        child: ElevatedButton(
          onPressed: () async {
            await signOut();
          },
          child: Text('Sign Out'),
        ),
      ),
    );
  }
}
