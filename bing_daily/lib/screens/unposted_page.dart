import 'package:bing_daily/constants.dart';
import 'package:flutter/material.dart';

/// Home page displaying a welcome message for the Bing Daily app.
class UnpostedPage extends StatelessWidget {
  const UnpostedPage({super.key});

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(title: const Text("Prompt of the Day:"), backgroundColor: bingGreen),
      body: const Center(
        child: Text(
          'Put image of Baxter here!!',
          style: TextStyle(color: bingGreen, fontSize: 20),
        ),
      ),
      backgroundColor: bingWhite,
    );
  }
}
