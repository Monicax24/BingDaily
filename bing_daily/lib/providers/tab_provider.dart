import 'package:flutter_riverpod/flutter_riverpod.dart';

/// Provider for the currently selected bottom navigation tab index.
final tabNotifierProvider = NotifierProvider<TabNotifier, int>(TabNotifier.new);

class TabNotifier extends Notifier<int> {
  @override
  int build() => 0;

  /// Sets the selected tab index.
  void setTab(int index) {
    state = index;
  }
}
