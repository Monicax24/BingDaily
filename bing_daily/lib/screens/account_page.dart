import 'package:bing_daily/api/firebase_auth.dart';
import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';

class AccountPage extends ConsumerWidget {
  const AccountPage({super.key});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    return Scaffold(
      appBar: AppBar(title: const Text('Bing Daily')),
      body: Column(
        mainAxisAlignment: MainAxisAlignment.center,
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          ElevatedButton(
            onPressed: () async {
              await signOut();
            },
            child: Text('Sign Out'),
          ),
          ElevatedButton(
            onPressed: () async {
            },
            child: Text('Print User Data'),
          ),
        ],
      ),
    );
  }
}
