from core.models import UserRepo, Issue
from core.serializers import UserRepoSerializer, IssueSerializer
from rest_framework import generics
import django_filters.rest_framework


class UserRepoList(generics.ListCreateAPIView):
    """
    Returns a list of issues, by optionally filtering against
    `repo` and `user` query parameter in the URL.
    """
    queryset = UserRepo.objects.all()
    serializer_class = UserRepoSerializer
    filter_backends = (django_filters.rest_framework.DjangoFilterBackend,)
    filter_fields = ('repo', 'user',)


class IssueList(generics.ListCreateAPIView):
    """
    Returns a list of issues, by optionally filtering against
    `language`, `tech_stack`, `experience_needed` and `expected_time`
    query parameter in the URL.
    """
    queryset = Issue.objects.all()
    serializer_class = IssueSerializer
    filter_backends = (django_filters.rest_framework.DjangoFilterBackend,)
    filter_fields = ('language', 'tech_stack', 'experience_needed', 'expected_time',)
