import 'package:bing_daily/constants.dart';
import 'package:bing_daily/screens/account_page.dart';
import 'package:bing_daily/screens/camera_page.dart';
import 'package:bing_daily/screens/home_page.dart';
import 'package:flutter/material.dart';

/// Main screen with bottom navigation bar for Home, Camera, and Account pages.
class MainScreen extends StatefulWidget {
  const MainScreen({super.key});

  @override
  State<MainScreen> createState() => _MainScreenState();
}

class _MainScreenState extends State<MainScreen> {
  int _selectedIndex = 0;

  // List of pages for navigation
  static const List<Widget> _pages = [HomePage(), CameraPage(), AccountPage()];

  /// Updates the selected page index when a nav item is tapped.
  void _onItemTapped(int index) {
    setState(() {
      _selectedIndex = index;
    });
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      body: _pages[_selectedIndex],
      bottomNavigationBar: BottomNavigationBar(
        items: const [
          BottomNavigationBarItem(
            icon: Icon(Icons.home_outlined),
            label: 'Home',
          ),
          BottomNavigationBarItem(
            icon: Icon(Icons.add_box_outlined),
            label: 'Post',
          ),
          BottomNavigationBarItem(
            icon: Icon(Icons.person_outline_rounded),
            label: 'Account',
          ),
        ],
        currentIndex: _selectedIndex,
        selectedItemColor: bingGreen,
        unselectedItemColor: Colors.grey,
        backgroundColor: bingWhite,
        onTap: _onItemTapped,
      ),
    );
  }
}
