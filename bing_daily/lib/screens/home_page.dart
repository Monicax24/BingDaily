import 'package:bing_daily/constants.dart';
import 'package:flutter/material.dart';

/// Home page displaying a welcome message for the Bing Daily app.
class HomePage extends StatefulWidget {
  const HomePage({super.key});

  @override
  State<HomePage> createState() => _HomePageState();
}

bool postedToday = false;

class _HomePageState extends State<HomePage> {
  @override
  Widget build(BuildContext context) {
    return Scaffold(
      body: Padding(
        padding: const EdgeInsets.only(top: 40.0),
        child: Column(
          mainAxisAlignment: MainAxisAlignment.center,
          crossAxisAlignment: CrossAxisAlignment.center,
          children: [
            const Padding(
              padding: EdgeInsets.symmetric(horizontal: 16.0),
              child: Text(
                'Hold up ...',
                style: TextStyle(
                  fontSize: 28,
                  fontWeight: FontWeight.bold,
                  color: Colors.black87,
                ),
                textAlign: TextAlign.center,
              ),
            ),
            SizedBox(height: 16),
            const Padding(
              padding: EdgeInsets.all(16.0),
              child: Text(
                'Post before you scroll',
                style: TextStyle(
                  fontSize: 18,
                  fontWeight: FontWeight.bold,
                  color: Colors.black87,
                ),
                textAlign: TextAlign.center,
              ),
            ),
            SizedBox(height: 16),
            ElevatedButton(
              onPressed: () {
                postedToday = true;
              },
              style: ElevatedButton.styleFrom(
                backgroundColor: bingGreen,
                padding: const EdgeInsets.symmetric(
                  horizontal: 30.0,
                  vertical: 5.0,
                ),
                shape: RoundedRectangleBorder(
                  borderRadius: BorderRadius.circular(8.0),
                ),
              ),
              child: const Text(
                'Post Now',
                style: TextStyle(
                  fontSize: 14,
                  fontWeight: FontWeight.bold,
                  color: Colors.white,
                ),
              ),
            ),
            SizedBox(height: 38),
            Container(
              height: 440,
              decoration: BoxDecoration(
                image: DecorationImage(
                  image: AssetImage('assets/images/baxter_stopsign.png'),
                  fit: BoxFit.fitHeight,
                ),
              ),
            ),
          ],
        ),
      ),
      backgroundColor: bingWhite,
    );
  }
}
