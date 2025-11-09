import 'package:bing_daily/constants.dart';
import 'package:bing_daily/providers/tab_provider.dart';
import 'package:bing_daily/screens/account_page.dart';
import 'package:bing_daily/screens/camera_page.dart';
import 'package:bing_daily/screens/home_page.dart';
import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';

/// Main screen with bottom navigation bar for Home, Camera, and Account pages.
class MainScreen extends ConsumerStatefulWidget {
  const MainScreen({super.key});

  @override
  ConsumerState<MainScreen> createState() => _MainScreenState();
}

class _MainScreenState extends ConsumerState<MainScreen> {
  // List of pages for navigation
  static const List<Widget> _pages = [HomePage(), CameraPage(), AccountPage()];

  /// Updates the selected page index when a nav item is tapped.
  void _onItemTapped(int index) {
    ref.read(tabNotifierProvider.notifier).setTab(index);
  }

  @override
  Widget build(BuildContext context) {
    final selectedIndex = ref.watch(tabNotifierProvider);

    return Scaffold(
      body: _pages[selectedIndex],
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
        currentIndex: selectedIndex,
        selectedItemColor: bingGreen,
        unselectedItemColor: Colors.grey,
        backgroundColor: bingWhite,
        onTap: _onItemTapped,
      ),
    );
  }
}
