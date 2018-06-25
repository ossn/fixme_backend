from core.models import UserRepo, Issue
from core.serializers import UserRepoSerializer, IssueSerializer
from rest_framework import generics
from  django_filters.rest_framework import DjangoFilterBackend
from rest_framework.filters import OrderingFilter
from rest_framework.views import APIView
from rest_framework.renderers import JSONRenderer, BrowsableAPIRenderer
from rest_framework.response import Response

class UserRepoList(generics.ListAPIView):
    """
    Returns a list of issues, by optionally filtering against
    `repo` and `user` query parameter in the URL.
    """
    queryset = UserRepo.objects.all()
    serializer_class = UserRepoSerializer
    filter_backends = (DjangoFilterBackend,)
    filter_fields = ('repo', 'user',)


class IssueList(generics.ListAPIView):
    """
    Returns a list of issues, by optionally filtering against
    `language`, `tech_stack`, `experience_needed` and `expected_time`
    query parameter in the URL.
    """
    queryset = Issue.objects.all()
    serializer_class = IssueSerializer
    filter_backends = (DjangoFilterBackend, OrderingFilter,)
    filter_fields = ('language', 'tech_stack', 'experience_needed', 'expected_time',)
    ordering_fields = ('experience_needed', 'expected_time')


class MetaData(APIView):
    """
    Returns a list of all the `language`, `tech_stack`,
    `experience_needed` present in the database.
    """
    renderer_classes = (JSONRenderer, BrowsableAPIRenderer,)

    def get(self, request, format=None):
        language = 'language'
        tech_stack = 'tech_stack'
        experience_needed = 'experience_needed'

        queryset = Issue.objects.values(
            language, tech_stack, experience_needed,)
        meta_data = {}
        meta_data[language] = set([tup[language] for tup in queryset])
        meta_data[tech_stack] = set([tup[tech_stack] for tup in queryset])
        meta_data[experience_needed] = set(
            [tup[experience_needed] for tup in queryset])
        return Response(meta_data)
