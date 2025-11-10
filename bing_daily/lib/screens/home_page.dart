import 'package:bing_daily/constants.dart';
import 'package:bing_daily/data/mock_data.dart';
import 'package:bing_daily/providers/has_posted_provider.dart';
import 'package:bing_daily/providers/tab_provider.dart';
import 'package:bing_daily/widgets/post_card.dart';
import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';

/// Home page that checks if user has posted today and displays posts or blocking UI.
class HomePage extends ConsumerWidget {
  const HomePage({super.key});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final hasPosted = ref.watch(hasPostedProvider);

    return Scaffold(
      body: Padding(
        padding: const EdgeInsets.only(top: 40.0),
        child: hasPosted ? _buildPostsView() : _buildBlockingView(context, ref),
      ),
      backgroundColor: bingWhite,
    );
  }

  /// Builds the view showing list of posts when user has posted.
  Widget _buildPostsView() {
    return ListView.builder(
      itemCount: mockPosts.length,
      itemBuilder: (context, index) => PostCard(post: mockPosts[index]),
    );
  }

  /// Builds the blocking view with prompt to post.
  Widget _buildBlockingView(BuildContext context, WidgetRef ref) {
    return Column(
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
        const SizedBox(height: 16),
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
        const SizedBox(height: 16),
        ElevatedButton(
          onPressed: () {
            ref.read(tabNotifierProvider.notifier).setTab(1);
            ref.read(hasPostedProvider.notifier).state = true;
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
        const SizedBox(height: 38),
        Container(
          height: 440,
          decoration: const BoxDecoration(
            image: DecorationImage(
              image: AssetImage('assets/images/baxter_stopsign.png'),
              fit: BoxFit.fitHeight,
            ),
          ),
        ),
      ],
    );
  }
}
